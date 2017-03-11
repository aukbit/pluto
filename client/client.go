package client

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/aukbit/pluto/client/balancer"
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// DefaultName prefix client name
	DefaultName    = "client"
	defaultVersion = "v1.0.0"
)

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type Client struct {
	cfg        *Config
	balancer   *balancer.Balancer    // Load Balancer to manage client connections
	requestsCh chan balancer.Request //
	connsCh    balancer.ConnsCh      //
	health     *health.Server        // Server implements `service Health`.
	logger     *zap.Logger           // client logger
}

// New create a new client
func New(opts ...Option) *Client {
	return newClient(opts...)
}

// newClient will instantiate a new Client with the given config
func newClient(opts ...Option) *Client {
	c := &Client{
		cfg:        newConfig(),
		balancer:   balancer.New(),
		requestsCh: make(chan balancer.Request),
		connsCh:    make(balancer.ConnsCh),
		health:     health.NewServer(),
	}
	c.logger, _ = zap.NewProduction()
	if len(opts) > 0 {
		c = c.WithOptions(opts...)
	}
	return c
}

// WithOptions clones the current Client, applies the supplied Options, and
// returns the resulting Client. It's safe to use concurrently.
func (c *Client) WithOptions(opts ...Option) *Client {
	d := c.clone()
	for _, opt := range opts {
		opt.apply(d)
	}
	return d
}

// clone creates a shallow copy client
func (c *Client) clone() *Client {
	copy := *c
	return &copy
}

func (c *Client) Config() *Config {
	return c.cfg
}

func (c *Client) initConnectors() error {
	if len(c.cfg.Targets) == 0 {
		return fmt.Errorf("connectors will not be initialized because targets were not provided")
	}
	for _, t := range c.cfg.Targets {

		conn := balancer.NewConnector(
			balancer.Target(t),
			balancer.GRPCRegister(c.cfg.GRPCRegister),
			balancer.UnaryClientInterceptors(c.cfg.UnaryClientInterceptors),
			balancer.Logger(c.logger),
		)
		// establish connection
		if err := conn.Init(); err != nil {
			return err
		}
		// add connector to the balancer pool
		c.balancer.Push(conn)
	}
	return nil
}

func (c *Client) startBalancer() {
	go c.balancer.Balance(c.requestsCh)
}

func (c *Client) Dial(opts ...Option) error {
	// set last configs
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(c)
		}
	}
	// register at service discovery
	if err := c.register(); err != nil {
		return err
	}
	// set targets from service discovery
	if err := c.targets(); err != nil {
		return err
	}
	// set logger
	c.logger = c.logger.With(
		zap.String("id", c.cfg.ID),
		zap.String("name", c.cfg.Name),
		zap.String("format", c.cfg.Format),
	)
	// init connectors
	if err := c.initConnectors(); err != nil {
		return err
	}
	// start load balancer
	c.startBalancer()
	// health check
	c.healthCheck()
	//
	return nil
}

// Request requests client for a connection to make a call
// func (c *Client) Request() balancer.Connector {
func (c *Client) Request() *balancer.Connector {
	r := balancer.NewRequest(c.connsCh)
	// send the request over the calls channel
	c.requestsCh <- r
	// return connector from connsCh
	return <-c.connsCh
}

// Done tells client that we no longer need the connection
// this methods is critical for load balancer to plays accordingly
// when service discovery is active
func (c *Client) Done(conn *balancer.Connector) {
	// tell balancer request it's Done
	c.balancer.Done(conn)
}

// Call acts as a shortcut. e.g if load balancer is serving a single service
// calling Request() and then Done(conn) is redundant so this could be
// wrapped internally and expose only the interface to be called directly on views
func (c *Client) Call() interface{} {
	conn := c.Request()
	defer c.Done(conn)
	return conn.Client()
}

func (c *Client) closeConnectors() {
	for _, n := range c.balancer.Pool() {
		n.Close()
	}
}

func (c *Client) Close() {
	c.logger.Info("close")
	// set health as not serving
	c.health.SetServingStatus(c.cfg.ID, 2)
	// close connectors
	c.closeConnectors()
	// unregister
	c.unregister()
}

// perform client health check on a random connector
func (c *Client) healthCheck() {
	// request a connector
	conn := c.Request()
	// connector health check
	if ok := conn.Health(); !ok {
		c.health.SetServingStatus(c.cfg.ID, 2)
		// TODO remove error connector and perform health check immediately!
		return
	}
	c.health.SetServingStatus(c.cfg.ID, 1)
}

// Health health check on client take in consideration
// a round trip to a server
func (c *Client) Health() *healthpb.HealthCheckResponse {
	// perform health check on a server
	c.healthCheck()
	// check client status
	hcr, err := c.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: c.cfg.ID})
	if err != nil {
		c.logger.Error("Health", zap.String("err", err.Error()))
		return &healthpb.HealthCheckResponse{Status: 2}
	}
	return hcr
}

// target get target IP:Port from service discovery system
func (c *Client) targets() error {
	if _, ok := c.cfg.Discovery.(discovery.Discovery); ok {
		t, err := c.cfg.Discovery.Service(c.cfg.TargetName)
		if err != nil {
			return err
		}
		c.cfg.Targets = t
	}
	return nil
}

// register Client within the service discovery system
func (c *Client) register() error {
	if _, ok := c.cfg.Discovery.(discovery.Discovery); ok {
		// define service
		dse := discovery.Service{
			Name: c.cfg.Name,
			Tags: []string{c.cfg.ID, c.cfg.Version},
		}
		// define check
		dck := discovery.Check{
			Name:  fmt.Sprintf("Service '%s' check", c.cfg.Name),
			Notes: fmt.Sprintf("Ensure the client is able to connect to service %s", c.cfg.TargetName),
			DeregisterCriticalServiceAfter: "10m",
			HTTP:      fmt.Sprintf("http://%s:9090/_health/client/%s", common.IPaddress(), c.cfg.Name),
			Interval:  "10s",
			Timeout:   "1s",
			ServiceID: c.cfg.Name,
		}
		if err := c.cfg.Discovery.Register(discovery.ServicesCfg(dse), discovery.ChecksCfg(dck)); err != nil {
			return err
		}
	}
	return nil
}

// unregister Server from the service discovery system
func (c *Client) unregister() error {
	if _, ok := c.cfg.Discovery.(discovery.Discovery); ok {
		if err := c.cfg.Discovery.Unregister(); err != nil {
			return err
		}
	}
	return nil
}

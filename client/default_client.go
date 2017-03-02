package client

import (
	"fmt"

	"github.com/aukbit/pluto/client/balancer"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type defaultClient struct {
	cfg        *Config
	balancer   *balancer.Balancer    // Load Balancer to manage client connections
	requestsCh chan balancer.Request //
	connsCh    balancer.ConnsCh      //
	health     *health.Server        // Server implements `service Health`.
	logger     zap.Logger            // client logger
}

// newClient will instantiate a new Client with the given config
func newClient(cfgs ...ConfigFn) *defaultClient {
	c := newConfig(cfgs...)
	return &defaultClient{
		cfg:        c,
		balancer:   balancer.NewBalancer(),
		requestsCh: make(chan balancer.Request),
		connsCh:    make(balancer.ConnsCh),
		health:     health.NewServer(),
		logger:     zap.New(zap.NewJSONEncoder())}
}

func (dc *defaultClient) Config() *Config {
	return dc.cfg
}

func (dc *defaultClient) initConnectors() error {
	if len(dc.cfg.Targets) == 0 {
		return fmt.Errorf("connectors will not be initialized because targets were not provided")
	}
	for _, t := range dc.cfg.Targets {
		c := balancer.NewConnector(
			balancer.Target(t),
			balancer.ParentID(dc.cfg.ID),
			balancer.GRPCRegister(dc.cfg.GRPCRegister),
			balancer.UnaryClientInterceptors(dc.cfg.UnaryClientInterceptors),
		)
		// establish connection
		if err := c.Init(); err != nil {
			return err
		}
		// add connector to the balancer pool
		dc.balancer.Push(c)
	}
	return nil
}

func (dc *defaultClient) startBalancer() {
	go dc.balancer.Balance(dc.requestsCh)
}

func (dc *defaultClient) Dial(cfgs ...ConfigFn) error {
	// set last configs
	for _, c := range cfgs {
		c(dc.cfg)
	}
	// register at service discovery
	if err := dc.register(); err != nil {
		return err
	}
	// set targets from service discovery
	if err := dc.targets(); err != nil {
		return err
	}
	// init logger
	dc.initLogger()
	//
	dc.logger.Info("start")
	// init connectors
	if err := dc.initConnectors(); err != nil {
		return err
	}
	// start load balancer
	dc.startBalancer()
	// health check
	dc.healthCheck()
	//
	return nil
}

// Request requests client for a connection to make a call
func (dc *defaultClient) Request() balancer.Connector {
	r := balancer.NewRequest(dc.connsCh)
	// send the request over the calls channel
	dc.requestsCh <- r
	// return connector from connsCh
	return <-dc.connsCh
}

// Done tells client that we no longer need the connection
// this methods is critical for load balancer to plays accordingly
// when service discovery is active
func (dc *defaultClient) Done(conn balancer.Connector) {
	// tell balancer request it's Done
	dc.balancer.Done(conn.Connector())
}

// Call acts as a shortcut. e.g if load balancer is serving a single service
// calling Request() and then Done(conn) is redundant so this could be
// wrapped internally and expose only the interface to be called directly on views
func (dc *defaultClient) Call() interface{} {
	conn := dc.Request()
	defer dc.Done(conn)
	return conn.Client()
}

func (dc *defaultClient) closeConnectors() {
	for _, c := range dc.balancer.Pool() {
		c.Close()
	}
}

func (dc *defaultClient) Close() {
	dc.logger.Info("close")
	// set health as not serving
	dc.health.SetServingStatus(dc.cfg.ID, 2)
	// close connectors
	dc.closeConnectors()
	// unregister
	dc.unregister()
}

// perform client health check on a random connector
func (dc *defaultClient) healthCheck() {
	// request a connector
	conn := dc.Request()
	// connector health check
	if ok := conn.Health(); !ok {
		dc.health.SetServingStatus(dc.cfg.ID, 2)
		// TODO remove error connector and perform health check immediately!
		return
	}
	dc.health.SetServingStatus(dc.cfg.ID, 1)
}

// Health health check on client take in consideration
// a round trip to a server
func (dc *defaultClient) Health() *healthpb.HealthCheckResponse {
	// perform health check on a server
	dc.healthCheck()
	// check client status
	hcr, err := dc.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: dc.cfg.ID})
	if err != nil {
		dc.logger.Error("Health", zap.String("err", err.Error()))
		return &healthpb.HealthCheckResponse{Status: 2}
	}
	return hcr
}

func (dc *defaultClient) initLogger() {
	dc.logger = dc.logger.With(
		zap.Nest("client",
			zap.String("id", dc.cfg.ID),
			zap.String("name", dc.cfg.Name),
			zap.String("format", dc.cfg.Format),
			zap.String("parent", dc.cfg.ParentID)))
}

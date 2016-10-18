package client

import (
	"errors"

	"bitbucket.org/aukbit/pluto/client/balancer"
	"github.com/uber-go/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type defaultClient struct {
	cfg    *Config
	logger zap.Logger

	call       interface{}           // TODO: deprecated managed in balancer
	conn       *grpc.ClientConn      // TODO: deprecated managed in balancer
	healthCall healthpb.HealthClient // TODO: deprecated managed in balancer
	health     *health.Server

	// Load Balancer to manage client connections
	balancer   *balancer.Balancer
	requestsCh chan balancer.Request
	connsCh    chan *balancer.Connector
}

// newClient will instantiate a new Client with the given config
func newClient(cfgs ...ConfigFn) *defaultClient {
	c := newConfig(cfgs...)
	return &defaultClient{
		cfg:        c,
		logger:     zap.New(zap.NewJSONEncoder()),
		health:     health.NewServer(),
		balancer:   balancer.NewBalancer(),
		requestsCh: make(chan balancer.Request),
		connsCh:    make(chan *balancer.Connector)}
}

func (dc *defaultClient) Config() *Config {
	return dc.cfg
}

func (dc *defaultClient) initConnectors() {
	for _, t := range dc.cfg.Targets {
		c := balancer.NewConnector(
			balancer.Target(t),
			balancer.GRPCRegister(dc.cfg.GRPCRegister),
			balancer.UnaryClientInterceptors(dc.cfg.UnaryClientInterceptors),
		)
		// establish connection
		if err := c.Dial(); err != nil {
			dc.logger.Error("dial", zap.String("err", err.Error()))
			return
		}
		// add connector to the balancer pool
		dc.balancer.Push(c)
		// watch for requests
		go c.Watch()
	}
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
	// if err := dc.register(); err != nil {
	// 	return err
	// }
	// set target from service discovery
	// if err := dc.target(); err != nil {
	// 	return err
	// }
	// init logger
	dc.initLogger()
	// init connectors
	dc.initConnectors()
	// start load balancer
	dc.startBalancer()
	// // start server
	// if err := dc.dialGRPC(); err != nil {
	// 	return err
	// }
	// set health
	dc.health.SetServingStatus(dc.cfg.ID, 1)
	//
	return nil
}

func (dc *defaultClient) Request() *balancer.Connector {
	r := balancer.NewRequest(dc.connsCh)
	// send the request over the calls channel
	dc.requestsCh <- r
	// return connector from connsCh
	return <-dc.connsCh
}

func (dc *defaultClient) Done(conn *balancer.Connector) {
	// send conn over balancer connsCh
	dc.balancer.ConnsCh <- conn
}

// TODO deprecated
func (dc *defaultClient) DialOld(cfgs ...ConfigFn) error {
	// set last configs
	for _, c := range cfgs {
		c(dc.cfg)
	}
	// register at service discovery
	if err := dc.register(); err != nil {
		return err
	}
	// set target from service discovery
	if err := dc.target(); err != nil {
		return err
	}
	// set logger
	dc.initLogger()
	// start server
	if err := dc.dialGRPC(); err != nil {
		return err
	}
	// set health
	dc.health.SetServingStatus(dc.cfg.ID, 1)
	//
	return nil
}

// TODO deprecated
func (dc *defaultClient) Call() interface{} {
	if dc.call == nil {
		return errors.New("Client has not been registered")
	}
	return dc.call
}

func (dc *defaultClient) Close() {
	dc.logger.Info("close")
	// set health as not serving
	dc.health.SetServingStatus(dc.cfg.ID, 2)
	// unregister
	// dc.unregister()
	// stop connectors
	for _, c := range dc.balancer.Pool() {
		go c.Stop()
	}
}

func (dc *defaultClient) healthServer() {
	// make health call on server
	hcr, err := dc.healthCall.Check(
		context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		dc.logger.Error("Health", zap.String("err", err.Error()))
		dc.health.SetServingStatus(dc.cfg.ID, 2)
		return
	}
	dc.health.SetServingStatus(dc.cfg.ID, hcr.Status)
}

// Health health check on client take in consideration
// the health check on server
func (dc *defaultClient) Health() *healthpb.HealthCheckResponse {
	//
	dc.healthServer()
	// make health call on client
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
			zap.String("target", dc.cfg.Target),
			zap.String("parent", dc.cfg.ParentID)))
}

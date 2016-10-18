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
	cfg        *Config
	logger     zap.Logger
	call       interface{}           // TODO: deprecated managed in balancer
	conn       *grpc.ClientConn      // TODO: deprecated managed in balancer
	healthCall healthpb.HealthClient // TODO: deprecated managed in balancer
	health     *health.Server

	// Load Balancer to manage client connections
	balancer *balancer.Balancer
	conns    []*balancer.Connector
}

// newClient will instantiate a new Client with the given config
func newClient(cfgs ...ConfigFn) *defaultClient {
	c := newConfig(cfgs...)
	return &defaultClient{
		cfg:      c,
		balancer: balancer.NewBalancer(),
		logger:   zap.New(zap.NewJSONEncoder()),
		health:   health.NewServer()}
}

func (dc *defaultClient) Config() *Config {
	return dc.cfg
}

func (dc *defaultClient) initConnectors() (conns []*balancer.Connector) {
	for _, t := range dc.cfg.Targets {
		c := balancer.NewConnector(
			balancer.Target(t),
			balancer.GRPCRegister(dc.cfg.GRPCRegister),
			balancer.UnaryClientInterceptors(dc.cfg.UnaryClientInterceptors),
		)
		conns = append(conns, c)
	}
	return conns
}

func (dc *defaultClient) initBalancer() error {
	return nil
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
	// init balancer
	dc.initBalancer()
	// start server
	if err := dc.dialGRPC(); err != nil {
		return err
	}
	// set health
	dc.health.SetServingStatus(dc.cfg.ID, 1)
	//
	return nil
}

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
	dc.unregister()
	// close connection
	dc.conn.Close()
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

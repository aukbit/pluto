package client

import (
	"errors"

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
	call       interface{}
	conn       *grpc.ClientConn
	healthCall healthpb.HealthClient
	health     *health.Server
}

// newClient will instantiate a new Client with the given config
func newClient(cfgs ...ConfigFunc) *defaultClient {
	c := newConfig(cfgs...)
	return &defaultClient{
		cfg:    c,
		logger: zap.New(zap.NewJSONEncoder()),
		health: health.NewServer()}
}

func (dc *defaultClient) Config() *Config {
	cfg := dc.cfg
	return cfg
}

func (dc *defaultClient) Dial(cfgs ...ConfigFunc) error {
	// set last configs
	for _, c := range cfgs {
		c(dc.cfg)
	}
	// set logger
	dc.setLogger()
	// start server
	if err := dc.dialGRPC(); err != nil {
		return err
	}
	// set health
	dc.health.SetServingStatus(dc.cfg.Name, 1)
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
	dc.health.SetServingStatus(dc.cfg.Name, 2)
	// close connection
	dc.conn.Close()
}

func (dc *defaultClient) healthServer() {
	// make health call on server
	hcr, err := dc.healthCall.Check(
		context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		dc.logger.Error("Health", zap.String("err", err.Error()))
		dc.health.SetServingStatus(dc.cfg.Name, 2)
		return
	}
	dc.health.SetServingStatus(dc.cfg.Name, hcr.Status)
}

// Health health check on client take in consideration
// the health check on server
func (dc *defaultClient) Health() *healthpb.HealthCheckResponse {
	//
	dc.healthServer()
	// make health call on client
	hcr, err := dc.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: dc.cfg.Name})
	if err != nil {
		dc.logger.Error("Health", zap.String("err", err.Error()))
		dc.health.SetServingStatus(dc.cfg.Name, 2)
		return &healthpb.HealthCheckResponse{Status: 2}
	}
	return hcr
}

func (dc *defaultClient) setLogger() {
	dc.logger = dc.logger.With(
		zap.Nest("client",
			zap.String("id", dc.cfg.ID),
			zap.String("name", dc.cfg.Name),
			zap.String("format", dc.cfg.Format),
			zap.String("target", dc.cfg.Target),
			zap.String("parent", dc.cfg.ParentID)))
}

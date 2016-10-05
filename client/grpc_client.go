package client

import (
	"errors"

	"github.com/uber-go/zap"

	"google.golang.org/grpc"
)

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type gRPCClient struct {
	cfg   *Config
	wire  interface{}
	close chan bool
}

// newGRPCClient will instantiate a new Client with the given config
func newClient(cfgs ...ConfigFunc) *gRPCClient {
	c := newConfig(cfgs...)
	return &gRPCClient{cfg: c, close: make(chan bool)}
}

func (g *gRPCClient) Init(cfgs ...ConfigFunc) error {
	for _, c := range cfgs {
		c(g.cfg)
	}
	return nil
}

func (g *gRPCClient) Config() *Config {
	cfg := g.cfg
	return cfg
}

func (g *gRPCClient) Dial() error {
	if err := g.dial(); err != nil {
		return err
	}
	return nil
}

func (g *gRPCClient) Call() interface{} {
	if g.wire == nil {
		return errors.New("gRPC client has not been registered")
	}
	return g.wire
}

func (g *gRPCClient) Close() error {
	// TODO
	g.close <- true
	return nil
}

func (g *gRPCClient) dial() error {
	logger.Info("dial",
		zap.String("client", g.cfg.Name),
		zap.String("id", g.cfg.ID),
		zap.String("format", g.cfg.Format),
	)
	// establishes gRPC client connection
	// TODO use TLS
	conn, err := grpc.Dial(g.Config().Target, grpc.WithInsecure())
	if err != nil {
		return err
	}
	// get gRPC client interface
	g.wire = g.cfg.RegisterClientFunc(conn)
	return nil
}

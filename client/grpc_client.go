package client

import (
	"errors"
	"log"

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
	log.Printf("DIAL  %s %s \t%s", g.cfg.Format, g.cfg.Name, g.cfg.ID)
	// establishes gRPC client connection
	// TODO use TLS
	conn, err := grpc.Dial(g.Config().Target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("ERROR %s grpc.Dial %v", g.cfg.Name, err)
	}
	// get gRPC client interface
	g.wire = g.cfg.RegisterClientFunc(conn)
	return nil
}

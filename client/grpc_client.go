package client

import (
	"log"
	"google.golang.org/grpc"
	"fmt"
)

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type gRPCClient struct {
	cfg 			*Config
	close 			chan bool
}

// newGRPCClient will instantiate a new Client with the given config
func newGRPCClient(cfgs ...ConfigFunc) Client {
	c := newConfig(cfgs...)
	c.Name = fmt.Sprintf("%s.grpc", c.Name)
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

func (g *gRPCClient) Dial() (i interface{}, err error) {
	if i, err = g.dial(); err != nil {
		return nil, err
	}
	return i, nil
}

func (g *gRPCClient) Close() error {
	// TODO
	g.close <-true
	return nil
}

func (g *gRPCClient) dial() (interface{}, error) {
	log.Printf("DIAL  %s %s", g.cfg.Name, g.cfg.Id)
	conn, err := grpc.Dial(g.Config().Target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("ERROR grpc.Dial %v", err)
	}
	return g.cfg.RegisterClientFunc(conn), nil
}
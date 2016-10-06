package client

import (
	"errors"

	"github.com/uber-go/zap"

	"google.golang.org/grpc"
)

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type gRPCClient struct {
	cfg          *Config
	logger       zap.Logger
	registration interface{}
	conn         *grpc.ClientConn
}

// newGRPCClient will instantiate a new Client with the given config
func newClient(cfgs ...ConfigFunc) *gRPCClient {
	c := newConfig(cfgs...)
	clt := &gRPCClient{cfg: c, logger: zap.New(zap.NewJSONEncoder())}
	return clt
}

func (g *gRPCClient) Config() *Config {
	cfg := g.cfg
	return cfg
}

func (g *gRPCClient) Dial(cfgs ...ConfigFunc) error {
	// set last configs
	for _, c := range cfgs {
		c(g.cfg)
	}
	// set logger
	g.setLogger()
	// start server
	if err := g.dial(); err != nil {
		return err
	}
	return nil
}

func (g *gRPCClient) Call() interface{} {
	if g.registration == nil {
		return errors.New("gRPC client has not been registered")
	}
	return g.registration
}

func (g *gRPCClient) Close() {
	g.logger.Info("close")
	g.conn.Close()
}

func (g *gRPCClient) dial() error {
	g.logger.Info("dial")
	// establishes gRPC client connection
	// TODO use TLS
	// append logger
	g.cfg.UnaryClientInterceptors = append(g.cfg.UnaryClientInterceptors, loggerUnaryClientInterceptor(g))
	// dial to establish connection
	conn, err := grpc.Dial(
		g.Config().Target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(WrapperUnaryClient(g.cfg.UnaryClientInterceptors...)))

	if err != nil {
		g.logger.Error("dial", zap.String("err", err.Error()))
		return err
	}
	// keep connection for later close
	g.conn = conn
	// register methods on connection
	g.registration = g.cfg.GRPCRegister(conn)
	return nil
}

func (g *gRPCClient) setLogger() {
	g.logger = g.logger.With(
		zap.Nest("client",
			zap.String("id", g.cfg.ID),
			zap.String("name", g.cfg.Name),
			zap.String("format", g.cfg.Format),
			zap.String("target", g.cfg.Target),
			zap.String("parent", g.cfg.ParentID)))
}

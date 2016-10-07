package client

import (
	"errors"

	"github.com/uber-go/zap"
	"google.golang.org/grpc"
)

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type defaultClient struct {
	cfg    *Config
	logger zap.Logger
	call   interface{}
	conn   *grpc.ClientConn
}

// newClient will instantiate a new Client with the given config
func newClient(cfgs ...ConfigFunc) *defaultClient {
	c := newConfig(cfgs...)
	return &defaultClient{cfg: c, logger: zap.New(zap.NewJSONEncoder())}
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
	dc.conn.Close()
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

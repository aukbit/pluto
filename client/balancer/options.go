package balancer

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Option is used to set options for the service.
type Option interface {
	apply(*Connector)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Connector)

func (f optionFunc) apply(s *Connector) {
	f(s)
}

// Target server address
func Target(t string) Option {
	return optionFunc(func(c *Connector) {
		c.cfg.Target = t
	})
}

// GRPCRegister register client gRPC function
func GRPCRegister(fn func(*grpc.ClientConn) interface{}) Option {
	return optionFunc(func(c *Connector) {
		c.cfg.GRPCRegister = fn
	})
}

// UnaryClientInterceptors ...
func UnaryClientInterceptors(uci []grpc.UnaryClientInterceptor) Option {
	return optionFunc(func(c *Connector) {
		c.cfg.UnaryClientInterceptors = append(c.cfg.UnaryClientInterceptors, uci...)
	})
}

// Logger sets a shallow copy from an input logger
func Logger(l *zap.Logger) Option {
	return optionFunc(func(c *Connector) {
		copy := *l
		c.logger = &copy
	})
}

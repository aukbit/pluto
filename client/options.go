package client

import (
	"time"

	"github.com/aukbit/pluto/common"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Option is used to set options for the service.
type Option interface {
	apply(*Client)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Client)

func (f optionFunc) apply(s *Client) {
	f(s)
}

// ID client id
func ID(id string) Option {
	return optionFunc(func(c *Client) {
		c.cfg.ID = id
	})
}

// Name client name
func Name(n string) Option {
	return optionFunc(func(c *Client) {
		c.cfg.Name = common.SafeName(n, DefaultName)
	})
}

// Description client description
func Description(d string) Option {
	return optionFunc(func(c *Client) {
		c.cfg.Description = d
	})
}

// Target server address
func Target(t string) Option {
	return optionFunc(func(c *Client) {
		c.cfg.Target = t
	})
}

// GRPCRegister register client gRPC function
func GRPCRegister(fn func(*grpc.ClientConn) interface{}) Option {
	return optionFunc(func(c *Client) {
		c.cfg.GRPCRegister = fn
	})
}

// UnaryClientInterceptors ...
func UnaryClientInterceptors(uci []grpc.UnaryClientInterceptor) Option {
	return optionFunc(func(c *Client) {
		c.cfg.UnaryClientInterceptors = append(c.cfg.UnaryClientInterceptors, uci...)
	})
}

// Logger sets a shallow copy from an input logger
func Logger(l zerolog.Logger) Option {
	return optionFunc(func(c *Client) {
		c.logger = l
	})
}

// Timeout returns an Option that configures a timeout for dialing a ClientConn
// initially. This is valid if and only if WithBlock() is present.
func Timeout(d time.Duration) Option {
	return optionFunc(func(c *Client) {
		c.cfg.Timeout = d
	})
}

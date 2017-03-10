package client

import (
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
	"github.com/aukbit/pluto/server"
	"go.uber.org/zap"
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

// Targets slice of server address
func Targets(t ...string) Option {
	return optionFunc(func(c *Client) {
		c.cfg.Targets = t
	})
}

// Target server address
func Target(t string) Option {
	return optionFunc(func(c *Client) {
		c.cfg.Targets[0] = t
	})
}

// TargetName server address
func TargetName(name string) Option {
	return optionFunc(func(c *Client) {
		c.cfg.TargetName = common.SafeName(name, server.DefaultName)
	})
}

// ParentID sets id of parent service that starts the server
func ParentID(id string) Option {
	return optionFunc(func(c *Client) {
		c.cfg.ParentID = id
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

// Discovery service discoery
func Discovery(d discovery.Discovery) Option {
	return optionFunc(func(c *Client) {
		c.cfg.Discovery = d
	})
}

// Development sets development logger
func Development() Option {
	return optionFunc(func(c *Client) {
		c.logger, _ = zap.NewDevelopment()
	})
}

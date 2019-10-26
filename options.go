package pluto

import (
	"github.com/aukbit/pluto/v6/client"
	"github.com/aukbit/pluto/v6/common"
	"github.com/aukbit/pluto/v6/discovery"
	"github.com/aukbit/pluto/v6/server"
	"github.com/rs/zerolog"
)

// Option is used to set options for the service.
type Option interface {
	apply(*Service)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Service)

func (f optionFunc) apply(s *Service) {
	f(s)
}

// ID service id
func ID(id string) Option {
	return optionFunc(func(s *Service) {
		s.cfg.ID = id
	})
}

// Name service name
func Name(n string) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Name = common.SafeName(n, defaultName)
	})
}

// Description service description
func Description(d string) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Description = d
	})
}

// Servers slice of service servers
func Servers(srv *server.Server) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Servers[srv.Name()] = srv
	})
}

// Clients slice of service clients
func Clients(clt *client.Client) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Clients[clt.Name()] = clt
		s.cfg.clientsCh <- clt
	})
}

// Discovery service discoery
func Discovery(d discovery.Discovery) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Discovery = d
	})
}

// HookAfterStart execute functions after service starts
func HookAfterStart(fn ...HookFunc) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Hooks["after_start"] = append(s.cfg.Hooks["after_start"], fn...)
	})
}

// Logger sets a new configuration for service logger
func Logger(l zerolog.Logger) Option {
	return optionFunc(func(s *Service) {
		s.logger = l
	})
}

// HealthAddr health server address
func HealthAddr(a string) Option {
	return optionFunc(func(s *Service) {
		s.cfg.HealthAddr = a
	})
}

package pluto

import (
	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/discovery"
	"github.com/aukbit/pluto/server"
	"go.uber.org/zap"
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
func Servers(srv server.Server) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Servers[srv.Config().Name] = srv
	})
}

// Clients slice of service clients
func Clients(clt client.Client) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Clients[clt.Config().Name] = clt
	})
}

// Datastore to persist data
func Datastore(d datastore.Datastore) Option {
	return optionFunc(func(s *Service) {
		s.cfg.Datastore = d
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

// Development sets development logger
func Development() Option {
	return optionFunc(func(s *Service) {
		s.logger, _ = zap.NewDevelopment()
	})
}

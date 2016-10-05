package pluto

import (
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/server"
)

// Service is the basic interface that defines what to expect from any server.
type Service interface {
	Init(...ConfigFunc) error
	Run() error
	Stop()
	Config() *Config
	Client(string) (client.Client, bool)
	Server(string) (server.Server, bool)
}

var (
	defaultName    = "pluto"
	defaultVersion = "1.0.0"
)

// NewService returns a new service with cfg passed in
func NewService(cfgs ...ConfigFunc) Service {
	return newService(cfgs...)
}

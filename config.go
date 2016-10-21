package pluto

import (
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/datastore"
	"bitbucket.org/aukbit/pluto/discovery"
	"bitbucket.org/aukbit/pluto/server"
	"github.com/uber-go/zap"
)

var logger = zap.New(zap.NewJSONEncoder())

// Config pluto service config
type Config struct {
	ID          string
	Name        string
	Description string
	Version     string
	Servers     map[string]server.Server
	Clients     map[string]client.Client
	Datastore   datastore.Datastore
	Discovery   discovery.Discovery
}

// ConfigFunc registers the Config
type ConfigFunc func(*Config)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Version: defaultVersion,
		Servers: make(map[string]server.Server),
		Clients: make(map[string]client.Client)}

	for _, c := range cfgs {
		c(cfg)
	}

	if len(cfg.ID) == 0 {
		cfg.ID = common.RandID("plt_", 6)
	}

	if len(cfg.Name) == 0 {
		cfg.Name = DefaultName
	}

	return cfg
}

// ID service id
func ID(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.ID = id
	}
}

// Name service name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Name = common.SafeName(n, DefaultName)
	}
}

// Description service description
func Description(d string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Description = d
	}
}

// Servers slice of service servers
func Servers(s server.Server) ConfigFunc {
	return func(cfg *Config) {
		cfg.Servers[s.Config().Name] = s
	}
}

// Clients slice of service clients
func Clients(c client.Client) ConfigFunc {
	return func(cfg *Config) {
		cfg.Clients[c.Config().Name] = c
	}
}

// Datastore to persist data
func Datastore(d datastore.Datastore) ConfigFunc {
	return func(cfg *Config) {
		cfg.Datastore = d
	}
}

// Discovery service discoery
func Discovery(d discovery.Discovery) ConfigFunc {
	return func(cfg *Config) {
		cfg.Discovery = d
	}
}

package pluto

import (
	"pluto/server"
	"fmt"
	"strings"
	"pluto/client"
)

type Config struct {
	Id 			string
	Name 			string
	Description 		string
	Version 		string
	Servers			[]server.Server
	Clients			[]client.Client
}

type ConfigFunc func(*Config)

var DefaultConfig = Config{
	Name: 			"default.pluto",
}

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := DefaultConfig

	for _, c := range cfgs {
		c(&cfg)
	}

	if len(cfg.Id) == 0 {
		cfg.Id = DefaultId
	}

	if len(cfg.Name) == 0 {
		cfg.Name = DefaultName
	}

	if len(cfg.Version) == 0 {
		cfg.Version = DefaultVersion
	}
	return &cfg
}

// Id service id
func Id(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Id = id
	}
}

// Name service name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Name = fmt.Sprintf("%s.%s", strings.Replace(strings.ToLower(n), " ", "", -1), DefaultName)
		for _, s := range cfg.Servers {
			s.Init(server.Name(cfg.Name))
		}

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
		s.Init(server.Name(cfg.Name))
		cfg.Servers = append(cfg.Servers, s)
	}
}

// Clients slice of service clients
func Clients(c client.Client) ConfigFunc {
	return func(cfg *Config) {
		c.Init(client.Name(cfg.Name))
		cfg.Clients = append(cfg.Clients, c)
	}
}
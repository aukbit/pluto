package pluto

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
	"pluto/server"
	"pluto/client"
)

type Config struct {
	Id 			string
	Name 			string
	Description 		string
	Version 		string
	Servers			map[string]server.Server
	Clients			map[string]client.Client
}

type ConfigFunc func(*Config)

var DefaultConfig = Config{
	Name: 			"default.pluto",
	Servers:		make(map[string]server.Server),
	Clients:		make(map[string]client.Client),
}

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := DefaultConfig

	for _, c := range cfgs {
		c(&cfg)
	}

	if len(cfg.Id) == 0 {
		cfg.Id = uuid.New().String()
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
		// TODO validate with regex
		cfg.Name = fmt.Sprintf("%s.%s", strings.Replace(strings.ToLower(n), " ", "", -1), DefaultName)
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
		//cfg.Servers = append(cfg.Servers, s)
		cfg.Servers[s.Config().Name] = s
	}
}

// Clients slice of service clients
func Clients(c client.Client) ConfigFunc {
	return func(cfg *Config) {
		//cfg.Clients = append(cfg.Clients, c)
		cfg.Clients[c.Config().Name] = c
	}
}
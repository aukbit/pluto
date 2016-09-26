package pluto

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
	"regexp"
	"log"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/datastore"
)

type Config struct {
	Id 			string
	Name 			string
	Description 		string
	Version 		string
	Servers			map[string]server.Server
	Clients			map[string]client.Client
	Datastore		datastore.Datastore
}

type ConfigFunc func(*Config)

var DefaultConfig = Config{
	Name: 			"pluto_default",
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
		// support only alphanumeric and underscore characters
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		safe := reg.ReplaceAllString(n, "_")
		cfg.Name = fmt.Sprintf("%s_%s", DefaultName, strings.ToLower(safe))
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
func Datastore(addr string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Datastore = datastore.NewDatastore(datastore.Addr(addr), datastore.Keyspace(cfg.Name))
	}
}

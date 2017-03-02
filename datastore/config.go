package datastore

import (
	"log"
	"regexp"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
)

type Config struct {
	ID         string
	Name       string
	Version    string
	Keyspace   string
	Target     string
	TargetName string // service name on service discovery
	Discovery  discovery.Discovery
}

type ConfigFunc func(*Config)

var (
	defaultKeyspace = "default"
	defaultTarget   = "127.0.0.1"
)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Version: defaultVersion,
		Keyspace: defaultKeyspace,
		Target:   defaultTarget}

	for _, c := range cfgs {
		c(cfg)
	}

	if len(cfg.ID) == 0 {
		cfg.ID = common.RandID("db_", 6)
	}

	if len(cfg.Name) == 0 {
		cfg.Name = DefaultName
	}

	return cfg
}

// ID client id
func ID(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.ID = id
	}
}

// Name client name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Name = common.SafeName(n, DefaultName)
	}
}

// Keyspace db keyspace
func Keyspace(ks string) ConfigFunc {
	return func(cfg *Config) {
		// cassandra valid characters
		//https://docs.datastax.com/en/cql/3.3/cql/cql_reference/ref-lexical-valid-chars.html
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		safe := reg.ReplaceAllString(ks, "_")
		cfg.Keyspace = safe
	}
}

// Target db address
func Target(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Target = a
	}
}

// TargetName server address
func TargetName(name string) ConfigFunc {
	return func(cfg *Config) {
		cfg.TargetName = name
	}
}

// Discovery service discoery
func Discovery(d discovery.Discovery) ConfigFunc {
	return func(cfg *Config) {
		cfg.Discovery = d
	}
}

package datastore

import (
	"log"
	"regexp"
)

type Config struct {
	Keyspace        string
	Addr            string
	TargetDiscovery string // service name on service discovery
}

type ConfigFunc func(*Config)

var (
	defaultKeyspace = "default"
	defaultAddr     = "127.0.0.1"
)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Keyspace: defaultKeyspace, Addr: defaultAddr}

	for _, c := range cfgs {
		c(cfg)
	}

	return cfg
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

// Addr db address
// TODO: rename Addr to Target
func Addr(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Addr = a
	}
}

// TargetDiscovery server address
func TargetDiscovery(name string) ConfigFunc {
	return func(cfg *Config) {
		cfg.TargetDiscovery = name
		// get target from service discovery
		// t, err := discovery.Target(name)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// cfg.Addr = t
	}
}

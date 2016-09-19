package datastore

import (
	"regexp"
	"log"
)

type Config struct {
	Keyspace		string
	Addr			string
}

type ConfigFunc func(*Config)

var DefaultConfig = Config{
	Keyspace:		"default",
	Addr:			"127.0.0.1",
}

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := DefaultConfig

	for _, c := range cfgs {
		c(&cfg)
	}

	return &cfg
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
func Addr(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Addr = a
	}
}


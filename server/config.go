package server

import (
	"fmt"
	"strings"
)


type Config struct {
	Id 			string
	Name 		string
	Description string
	Version 	string
	Addr       	string        // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
}

type ConfigFunc func(*Config)

var DefaultConfig = Config{
	Name: 			"default.server",
	Addr:			":8080",
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

// Server Id
func Id(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Id = id
	}
}

// Server name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Name = fmt.Sprintf("%s.%s", strings.ToLower(n), DefaultName)
	}
}

// Server description
func Description(d string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Description = d
	}
}

// Server description
func Addr(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Addr = a
	}
}
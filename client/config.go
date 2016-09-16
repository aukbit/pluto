package client

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Config struct {
	Id 			string
	Name 			string
	Description 		string
	Version 		string
	Target       		string        // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	RegisterClientFunc	func(*grpc.ClientConn) interface{}
}

type ConfigFunc func(*Config)

var DefaultConfig = Config{
	Name: 			"default.client",
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

// Id cleint id
func Id(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Id = id
	}
}

// Name client name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Name = fmt.Sprintf("%s.%s", strings.ToLower(n), DefaultName)
	}
}

// Description client description
func Description(d string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Description = d
	}
}

// Target server address
func Target(t string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Target = t
	}
}

// RegisterClientFunc register client gRPC function
func RegisterClientFunc(fn func(*grpc.ClientConn) interface{}) ConfigFunc {
	return func(cfg *Config) {
		cfg.RegisterClientFunc = fn
	}
}
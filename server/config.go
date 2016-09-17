package server

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
	"pluto/server/router"
	"google.golang.org/grpc"
)


type Config struct {
	Id			string
	Name 			string
	Description 		string
	Version			string
	Addr       		string        // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	Format			string
	Router			*router.Router
	RegisterServerFunc	func(*grpc.Server)
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

// Id server id
func Id(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Id = id
	}
}

// Name server name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Name = fmt.Sprintf("%s.%s", strings.ToLower(n), DefaultName)
	}
}

// Description server description
func Description(d string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Description = d
	}
}

// Addr server address
func Addr(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Addr = a
	}
}

// Router to be used on http servers
func Router(r *router.Router) ConfigFunc {
	return func(cfg *Config) {
		cfg.Router = r
	}
}

// RegisterServerFunc register gRPC server function
func RegisterServerFunc(fn func(*grpc.Server)) ConfigFunc {
	return func(cfg *Config) {
		cfg.RegisterServerFunc = fn
	}
}


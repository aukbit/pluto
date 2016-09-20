package server

import (
	"fmt"
	"strings"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"regexp"
	"bitbucket.org/aukbit/pluto/server/router"
)


type Config struct {
	Id			string
	Name 			string
	Description 		string
	Version			string
	Addr       		string        // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	Format			string
	Mux			router.Mux
	RegisterServerFunc	func(*grpc.Server)
}

type ConfigFunc func(*Config)


var DefaultConfig = Config{
	Name: 			"server_default",
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
		// support only alphanumeric and underscore characters
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		safe := reg.ReplaceAllString(n, "_")
		cfg.Name = fmt.Sprintf("%s_%s", DefaultName, strings.ToLower(safe))
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

// Mux server multiplexer
func Mux(m router.Mux) ConfigFunc {
	return func(cfg *Config) {
		cfg.Mux = m
	}
}

// RegisterServerFunc register gRPC server function
func RegisterServerFunc(fn func(*grpc.Server)) ConfigFunc {
	return func(cfg *Config) {
		cfg.RegisterServerFunc = fn
	}
}

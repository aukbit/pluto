package client

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Config struct {
	Id                 string
	Name               string
	Description        string
	Version            string
	Target             string // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	Format             string
	RegisterClientFunc func(*grpc.ClientConn) interface{}
}

type ConfigFunc func(*Config)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{}

	for _, c := range cfgs {
		c(cfg)
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
	return cfg
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
		// support only alphanumeric and underscore characters
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		safe := reg.ReplaceAllString(n, "_")
		cfg.Name = fmt.Sprintf("%s_%s", DefaultName, strings.ToLower(safe))
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

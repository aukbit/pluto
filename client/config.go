package client

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/discovery"
	"bitbucket.org/aukbit/pluto/server"

	"google.golang.org/grpc"
)

// Config client configuaration options
type Config struct {
	ID                      string
	Name                    string
	Description             string
	Version                 string
	Target                  string // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	TargetName              string // service name on service discovery
	Format                  string
	ParentID                string // sets parent ID
	GRPCRegister            func(*grpc.ClientConn) interface{}
	UnaryClientInterceptors []grpc.UnaryClientInterceptor // gRPC interceptors
	Discovery               discovery.Discovery
}

// ConfigFunc registers the Config
type ConfigFunc func(*Config)

var (
	defaultTarget = "localhost:65060"
	defaultFormat = "grpc"
)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Target: defaultTarget, Format: defaultFormat, Version: defaultVersion}

	for _, c := range cfgs {
		c(cfg)
	}

	if len(cfg.ID) == 0 {
		cfg.ID = common.RandID("clt_", 6)
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

// TargetName server address
func TargetName(name string) ConfigFunc {
	return func(cfg *Config) {
		// support only alphanumeric and underscore characters
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		safe := reg.ReplaceAllString(name, "_")
		cfg.TargetName = fmt.Sprintf("%s_%s", server.DefaultName, strings.ToLower(safe))
	}
}

// ParentID sets id of parent service that starts the server
func ParentID(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.ParentID = id
	}
}

// GRPCRegister register client gRPC function
func GRPCRegister(fn func(*grpc.ClientConn) interface{}) ConfigFunc {
	return func(cfg *Config) {
		cfg.GRPCRegister = fn
	}
}

// Discovery service discoery
func Discovery(d discovery.Discovery) ConfigFunc {
	return func(cfg *Config) {
		cfg.Discovery = d
	}
}

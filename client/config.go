package client

import (
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
	"github.com/aukbit/pluto/server"

	"google.golang.org/grpc"
)

// Config client configuaration options
type Config struct {
	ID          string
	Name        string
	Description string
	Version     string
	Targets     []string // slice of TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	// Target                  string   // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	TargetName              string // service name on service discovery
	Format                  string
	ParentID                string // sets parent ID
	GRPCRegister            func(*grpc.ClientConn) interface{}
	UnaryClientInterceptors []grpc.UnaryClientInterceptor // gRPC interceptors
	Discovery               discovery.Discovery
}

// ConfigFn registers the Config
type ConfigFn func(*Config)

var (
	defaultTarget = "localhost:65060"
	defaultFormat = "grpc"
)

func newConfig(cfgs ...ConfigFn) *Config {

	cfg := &Config{
		Targets: []string{defaultTarget},
		Format:  defaultFormat,
		Version: defaultVersion}

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

// Target return client Target address
func (c *Config) Target() string {
	if len(c.Targets) > 0 {
		return c.Targets[0]
	}
	return ""
}

// ID client id
func ID(id string) ConfigFn {
	return func(cfg *Config) {
		cfg.ID = id
	}
}

// Name client name
func Name(n string) ConfigFn {
	return func(cfg *Config) {
		cfg.Name = common.SafeName(n, DefaultName)
	}
}

// Description client description
func Description(d string) ConfigFn {
	return func(cfg *Config) {
		cfg.Description = d
	}
}

// Targets slice of server address
func Targets(t ...string) ConfigFn {
	return func(cfg *Config) {
		cfg.Targets = t
	}
}

// Target server address
func Target(t string) ConfigFn {
	return func(cfg *Config) {
		cfg.Targets[0] = t
	}
}

// TargetName server address
func TargetName(name string) ConfigFn {
	return func(cfg *Config) {
		cfg.TargetName = common.SafeName(name, server.DefaultName)
	}
}

// ParentID sets id of parent service that starts the server
func ParentID(id string) ConfigFn {
	return func(cfg *Config) {
		cfg.ParentID = id
	}
}

// GRPCRegister register client gRPC function
func GRPCRegister(fn func(*grpc.ClientConn) interface{}) ConfigFn {
	return func(cfg *Config) {
		cfg.GRPCRegister = fn
	}
}

// UnaryClientInterceptors ...
func UnaryClientInterceptors(uci []grpc.UnaryClientInterceptor) ConfigFn {
	return func(cfg *Config) {
		cfg.UnaryClientInterceptors = append(cfg.UnaryClientInterceptors, uci...)
	}
}

// Discovery service discoery
func Discovery(d discovery.Discovery) ConfigFn {
	return func(cfg *Config) {
		cfg.Discovery = d
	}
}

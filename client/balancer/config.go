package balancer

import (
	"bitbucket.org/aukbit/pluto/common"
	"google.golang.org/grpc"
)

// Config client configuaration options
type Config struct {
	ID                      string
	Name                    string
	Target                  string                             // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	ParentID                string                             // sets parent ID
	GRPCRegister            func(*grpc.ClientConn) interface{} //
	UnaryClientInterceptors []grpc.UnaryClientInterceptor      // gRPC interceptors
}

// ConfigFn registers the Config
type ConfigFn func(*Config)

var (
	defaultName   = "connector"
	defaultTarget = "localhost:65060"
)

func newConfig(cfgs ...ConfigFn) *Config {

	cfg := &Config{Target: defaultTarget}

	for _, c := range cfgs {
		c(cfg)
	}

	if len(cfg.ID) == 0 {
		cfg.ID = common.RandID("con_", 6)
	}

	if len(cfg.Name) == 0 {
		cfg.Name = defaultName
	}

	return cfg
}

// Target server address
func Target(t string) ConfigFn {
	return func(cfg *Config) {
		cfg.Target = t
	}
}

// ParentID sets id of parent client
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

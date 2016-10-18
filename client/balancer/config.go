package balancer

import "google.golang.org/grpc"

// Config client configuaration options
type Config struct {
	Target       string // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	GRPCRegister GRPCRegisterFn
	// UnaryClientInterceptors []grpc.UnaryClientInterceptor // gRPC interceptors
}

// ConfigFn registers the Config
type ConfigFn func(*Config)

// GRPCRegisterFn func type
type GRPCRegisterFn func(*grpc.ClientConn) interface{}

var (
	defaultTarget = "localhost:65060"
	defaultFormat = "grpc"
)

func newConfig(cfgs ...ConfigFn) *Config {

	cfg := &Config{Target: defaultTarget}

	for _, c := range cfgs {
		c(cfg)
	}

	return cfg
}

// Target server address
func Target(t string) ConfigFn {
	return func(cfg *Config) {
		cfg.Target = t
	}
}

// GRPCRegister register client gRPC function
func GRPCRegister(fn GRPCRegisterFn) ConfigFn {
	return func(cfg *Config) {
		cfg.GRPCRegister = fn
	}
}

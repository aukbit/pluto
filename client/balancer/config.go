package balancer

import "google.golang.org/grpc"

// Config client configuaration options
type Config struct {
	Target                  string                             // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	GRPCRegister            func(*grpc.ClientConn) interface{} //
	UnaryClientInterceptors []grpc.UnaryClientInterceptor      // gRPC interceptors
}

// ConfigFn registers the Config
type ConfigFn func(*Config)

var (
	defaultTarget = "localhost:65060"
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

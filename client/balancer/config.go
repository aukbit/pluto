package balancer

import (
	"google.golang.org/grpc"
)

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

func newConfig() *Config {
	return &Config{
		Target: defaultTarget,
	}
}

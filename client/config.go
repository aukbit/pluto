package client

import (
	"time"

	"github.com/aukbit/pluto/common"

	"google.golang.org/grpc"
)

// Config client configuaration options
type Config struct {
	ID                      string
	Name                    string
	Description             string
	Target                  string // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	Format                  string
	GRPCRegister            func(*grpc.ClientConn) interface{}
	UnaryClientInterceptors []grpc.UnaryClientInterceptor // gRPC interceptors
	Timeout                 time.Duration
}

var (
	defaultTarget = "localhost:65060"
	defaultFormat = "grpc"
)

func newConfig() *Config {
	return &Config{
		ID:      common.RandID("clt_", 6),
		Name:    DefaultName,
		Format:  defaultFormat,
		Timeout: 500 * time.Millisecond,
	}
}

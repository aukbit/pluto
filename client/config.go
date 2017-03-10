package client

import (
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"

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

var (
	defaultTarget = "localhost:65060"
	defaultFormat = "grpc"
)

func newConfig() *Config {
	return &Config{
		ID:      common.RandID("clt_", 6),
		Name:    DefaultName,
		Targets: []string{defaultTarget},
		Format:  defaultFormat,
		Version: defaultVersion,
	}
}

// Target return client Target address
func (c *Config) Target() string {
	if len(c.Targets) > 0 {
		return c.Targets[0]
	}
	return ""
}

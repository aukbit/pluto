package server

import (
	"crypto/tls"
	"log"
	"regexp"
	"strconv"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
	"github.com/aukbit/pluto/server/router"
	"google.golang.org/grpc"
)

// Config server configuaration options
type Config struct {
	ID                       string
	Name                     string
	Description              string
	Addr                     string // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	Format                   string
	Mux                      *router.Router
	TLSConfig                *tls.Config // optional TLS config, used by ListenAndServeTLS
	GRPCRegister             GRPCRegisterServiceFunc
	Middlewares              []router.Middleware            // http middlewares
	UnaryServerInterceptors  []grpc.UnaryServerInterceptor  // gRPC interceptors
	StreamServerInterceptors []grpc.StreamServerInterceptor // gRPC interceptors
	Discovery                discovery.Discovery
}

// GRPCRegisterServiceFunc grpc
type GRPCRegisterServiceFunc func(*grpc.Server)

var (
	defaultAddr   = ":8080"
	defaultFormat = "http"
)

func newConfig() *Config {
	return &Config{
		ID:     common.RandID("srv_", 6),
		Name:   DefaultName,
		Addr:   defaultAddr,
		Format: defaultFormat,
	}
}

// Convert string Addr to int Port
func (c *Config) Port() int {
	// support only numeric
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	safe := reg.ReplaceAllString(c.Addr, "")
	i, err := strconv.Atoi(safe)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

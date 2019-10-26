package server

import (
	"crypto/tls"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/aukbit/pluto/v6/common"
	"github.com/aukbit/pluto/v6/discovery"
	"github.com/aukbit/pluto/v6/server/router"
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
	Discovery                discovery.Discovery
	ReadTimeout              time.Duration
	WriteTimeout             time.Duration
	mu                       sync.Mutex                     // ensures atomic writes; protects the following fields
	Middlewares              []router.Middleware            // http middlewares
	UnaryServerInterceptors  []grpc.UnaryServerInterceptor  // gRPC interceptors
	StreamServerInterceptors []grpc.StreamServerInterceptor // gRPC interceptors
}

// GRPCRegisterServiceFunc grpc
type GRPCRegisterServiceFunc func(*grpc.Server)

var (
	defaultAddr   = ":8080"
	defaultFormat = "http"
)

func newConfig() Config {
	return Config{
		ID:           common.RandID("srv_", 6),
		Name:         DefaultName,
		Addr:         defaultAddr,
		Format:       defaultFormat,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// Port converts string Addr to int Port
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

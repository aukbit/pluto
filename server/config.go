package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"regexp"
	"strings"

	"bitbucket.org/aukbit/pluto/server/router"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// Config server configuaration options
type Config struct {
	ID           string
	Name         string
	Description  string
	Version      string
	Addr         string // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
	Format       string
	Mux          router.Mux
	TLSConfig    *tls.Config // optional TLS config, used by ListenAndServeTLS
	GRPCRegister GRPCRegisterServiceFunc
}

// GRPCRegisterServiceFunc grpc
type GRPCRegisterServiceFunc func(*grpc.Server)

// ConfigFunc registers the Config
type ConfigFunc func(*Config)

var (
	defaultAddr   = ":8080"
	defaultFormat = "http"
)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Addr: defaultAddr, Format: defaultFormat, Version: defaultVersion}

	for _, c := range cfgs {
		c(cfg)
	}

	if len(cfg.ID) == 0 {
		cfg.ID = uuid.New().String()
	}

	if len(cfg.Name) == 0 {
		cfg.Name = defaultName
	}

	return cfg
}

// ID server id
func ID(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.ID = id
	}
}

// Name server name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		// support only alphanumeric and underscore characters
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		safe := reg.ReplaceAllString(n, "_")
		cfg.Name = fmt.Sprintf("%s_%s", defaultName, strings.ToLower(safe))
	}
}

// Description server description
func Description(d string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Description = d
	}
}

// Addr server address
func Addr(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Addr = a
	}
}

// Mux server multiplexer
func Mux(m router.Mux) ConfigFunc {
	return func(cfg *Config) {
		cfg.Mux = m
	}
}

// TLSConfig server multiplexer
func TLSConfig(certFile, keyFile string) ConfigFunc {
	return func(cfg *Config) {
		cer, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Printf("ERROR tls.LoadX509KeyPair %v", err)
			return
		}
		cfg.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			Certificates:             []tls.Certificate{cer},
		}
		cfg.Format = "https"
	}
}

// RegisterClientFunc register client gRPC function
func GRPCRegister(fn GRPCRegisterServiceFunc) ConfigFunc {
	// func GRPCRegister(i interface{}) ConfigFunc {
	return func(cfg *Config) {
		cfg.GRPCRegister = fn
		cfg.Format = "grpc"
	}
}

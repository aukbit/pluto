package server

import (
	"crypto/tls"
	"log"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
	"github.com/aukbit/pluto/server/router"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

// Option is used to set options for the server.
type Option interface {
	apply(*Server)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Server)

func (f optionFunc) apply(s *Server) {
	f(s)
}

// ID server id
func ID(id string) Option {
	return optionFunc(func(s *Server) {
		s.cfg.ID = id
	})
}

// Name server name
func Name(n string) Option {
	return optionFunc(func(s *Server) {
		s.cfg.Name = common.SafeName(n, DefaultName)
	})
}

// Description server description
func Description(d string) Option {
	return optionFunc(func(s *Server) {
		s.cfg.Description = d
	})
}

// Addr server address
func Addr(a string) Option {
	return optionFunc(func(s *Server) {
		s.cfg.Addr = a
	})
}

// Mux server multiplexer
func Mux(m *router.Router) Option {
	return optionFunc(func(s *Server) {
		s.cfg.Mux = m
	})
}

// TLSConfig server multiplexer
func TLSConfig(certFile, keyFile string) Option {
	return optionFunc(func(s *Server) {
		cer, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Printf("ERROR tls.LoadX509KeyPair %v", err)
			return
		}
		s.cfg.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			Certificates:             []tls.Certificate{cer},
		}
		s.cfg.Format = "https"
	})
}

// GRPCRegister register client gRPC function
func GRPCRegister(fn GRPCRegisterServiceFunc) Option {
	return optionFunc(func(s *Server) {
		s.cfg.GRPCRegister = fn
		s.cfg.Format = "grpc"
	})
}

// Middlewares slice with router.Middleware
func Middlewares(m ...router.Middleware) Option {
	return optionFunc(func(s *Server) {
		s.cfg.Middlewares = append(s.cfg.Middlewares, m...)
	})
}

// UnaryServerInterceptors slice with grpc.UnaryServerInterceptor
func UnaryServerInterceptors(i ...grpc.UnaryServerInterceptor) Option {
	return optionFunc(func(s *Server) {
		s.cfg.UnaryServerInterceptors = append(s.cfg.UnaryServerInterceptors, i...)
	})
}

// StreamServerInterceptors slice with grpc.StreamServerInterceptor
func StreamServerInterceptors(i ...grpc.StreamServerInterceptor) Option {
	return optionFunc(func(s *Server) {
		s.cfg.StreamServerInterceptors = append(s.cfg.StreamServerInterceptors, i...)
	})
}

// Discovery service discoery
func Discovery(d discovery.Discovery) Option {
	return optionFunc(func(s *Server) {
		s.cfg.Discovery = d
	})
}

// Logger sets a shallow copy from an input logger
func Logger(l zerolog.Logger) Option {
	return optionFunc(func(s *Server) {
		s.logger = l
	})
}

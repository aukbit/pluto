package server

import (
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Server is the basic interface that defines what to expect from any server.
type Server interface {
	Run(...ConfigFunc) error
	Stop()
	Config() *Config
	Health() *healthpb.HealthCheckResponse
}

const (
	// DefaultName server prefix name
	DefaultName    = "plt_server"
	defaultVersion = "v1.0.0"
)

// NewServer returns a new http server with cfg passed in
func NewServer(cfgs ...ConfigFunc) Server {
	return newServer(cfgs...)
}

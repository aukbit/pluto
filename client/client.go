package client

import (
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Client is an interface to make calls to services
type Client interface {
	Dial(...ConfigFunc) error
	Call() interface{}
	Close()
	Config() *Config
	Health() *healthpb.HealthCheckResponse
}

const (
	// DefaultName prefix client name
	DefaultName    = "plt_client"
	defaultVersion = "v1.0.0"
)

// NewClient returns a new client with cfg passed in
func NewClient(cfgs ...ConfigFunc) Client {
	return newClient(cfgs...)
}

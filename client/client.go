package client

import (
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Client is an interface to make calls to services
type Client interface {
	Dial(...ConfigFunc) error
	Call() interface{}
	Health() healthpb.HealthClient
	Close()
	Config() *Config
}

var (
	defaultName    = "client"
	defaultVersion = "v1.0.0"
)

// NewClient returns a new client with cfg passed in
func NewClient(cfgs ...ConfigFunc) Client {
	return newClient(cfgs...)
}

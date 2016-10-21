package client

import (
	"bitbucket.org/aukbit/pluto/client/balancer"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Client is an interface to make calls to services
type Client interface {
	Dial(...ConfigFn) error
	Request() balancer.Connector
	Done(balancer.Connector)
	Call() interface{}
	Close()
	Config() *Config
	Health() *healthpb.HealthCheckResponse
}

const (
	// DefaultName prefix client name
	DefaultName    = "client"
	defaultVersion = "v1.0.0"
)

// NewClient returns a new client with cfg passed in
func NewClient(cfgs ...ConfigFn) Client {
	return newClient(cfgs...)
}

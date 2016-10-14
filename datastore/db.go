package datastore

import (
	"github.com/gocql/gocql"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Datastore ...
type Datastore interface {
	Connect(cfgs ...ConfigFunc) error
	Session() *gocql.Session
	RefreshSession() error
	Close()
	Config() *Config
	Health() *healthpb.HealthCheckResponse
}

const (
	// DefaultName prefix datastore client name
	DefaultName    = "plt_client_db"
	defaultVersion = "v1.0.0"
)

// NewDatastore ...
func NewDatastore(cfgs ...ConfigFunc) Datastore {
	return newDatastore(cfgs...)
}

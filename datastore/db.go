package datastore

import "github.com/gocql/gocql"

type Datastore interface {
	Connect(cfgs ...ConfigFunc)
	Session() *gocql.Session
	RefreshSession() error
	Close()
}

const (
	// DefaultName prefix datastore client name
	DefaultName    = "db_client"
	defaultVersion = "v1.0.0"
)

func NewDatastore(cfgs ...ConfigFunc) Datastore {
	return newDatastore(cfgs...)
}

package db

import "github.com/gocql/gocql"

type Datastore interface {
	Session()		*gocql.Session
	RefreshSession()
	Close()
}

func NewDatastore(cfgs ...ConfigFunc) Datastore {
	return newDatastore(cfgs...)
}
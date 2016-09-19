package datastore

import "github.com/gocql/gocql"

type Datastore interface {
	Connect()
	Session()		*gocql.Session
	RefreshSession()	error
	Close()
}

func NewDatastore(cfgs ...ConfigFunc) Datastore {
	return newDatastore(cfgs...)
}
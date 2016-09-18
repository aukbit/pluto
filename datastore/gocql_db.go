package datastore

import (
	"github.com/gocql/gocql"
	"log"
)

var (
	DefaultCluster="127.0.0.1"
)

type datastore struct {
	cfg 			*Config
	cluster			*gocql.ClusterConfig
	session			*gocql.Session
}

// NewServer will instantiate a new Server with the given config
func newDatastore(cfgs ...ConfigFunc) Datastore {
	c := newConfig(cfgs...)
	return &datastore{cfg: c}
}

func (ds *datastore) Connect() {
	log.Printf("----- %s cluster connected on %s", ds.cfg.Keyspace, ds.cfg.Addr)
	ds.cluster = gocql.NewCluster(ds.cfg.Addr)
	ds.cluster.Keyspace = ds.cfg.Keyspace
}

func (ds *datastore) RefreshSession() error {
	s, err := ds.cluster.CreateSession()
	if err != nil {
		return err
	}
	ds.session = s
	return nil
}

func (ds *datastore) Close() {
	ds.session.Close()
}

func (ds *datastore) Session() *gocql.Session {
	return ds.session
}

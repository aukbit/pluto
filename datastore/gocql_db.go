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
	log.Printf("----- gcql %s connected on %s", ds.cfg.Keyspace, ds.cfg.Addr)
	ds.cluster = gocql.NewCluster(ds.cfg.Addr)
	ds.cluster.ProtoVersion = 3
	ds.cluster.Keyspace = ds.cfg.Keyspace
}

func (ds *datastore) RefreshSession() error {
	s, err := ds.cluster.CreateSession()
	if err != nil {
		return err
		if err == gocql.ErrNoConnectionsStarted {
			// TODO currently gocql driver as an func createKeyspace
			return err
		}
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

func (ds *datastore) createKeyspace(keyspace string, replicationFactor int) error {
	q := "CREATE KEYSPACE ? WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : ? };"
	if err := ds.session.Query(q, keyspace, replicationFactor).Exec(); err != nil {
		return err
	}
	return nil
}

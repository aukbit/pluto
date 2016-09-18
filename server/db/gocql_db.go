package db

import (
	"github.com/gocql/gocql"
	"log"
)

var (
	DefaultCluster="127.0.0.1"
)

type datastore struct {
	cluster		*gocql.ClusterConfig
	session		*gocql.Session
}

// NewServer will instantiate a new Server with the given config
func newDatastore(cfgs ...ConfigFunc) Datastore {
	c := newConfig(cfgs...)
	cluster := gocql.NewCluster(c.Addr)
	cluster.Keyspace = c.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("ERROR cluster.CreateSession() %v", err)
	}
	return &datastore{cluster: cluster, session: session}
}


func (ds *datastore) RefreshSession() {
	ds.session = ds.cluster.CreateSession()
}

func (ds *datastore) Close() {
	ds.session.Close()
}

func (ds *datastore) Session() *gocql.Session {
	return ds.session
}

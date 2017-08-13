package datastore

import (
	"github.com/aukbit/pluto/common"
	"github.com/gocql/gocql"
	mgo "gopkg.in/mgo.v2"
)

type Config struct {
	ID        string
	Name      string
	driver    string
	Cassandra *gocql.ClusterConfig
	MongoDB   *mgo.DialInfo
}

func newConfig() Config {
	return Config{
		ID:   common.RandID("db_", 6),
		Name: defaultName,
	}
}

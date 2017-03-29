package datastore

import (
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
	"github.com/gocql/gocql"
	mgo "gopkg.in/mgo.v2"
)

type Config struct {
	ID         string
	Name       string
	Version    string
	Keyspace   string
	Target     string
	TargetName string // service name on service discovery
	Discovery  discovery.Discovery
	driver     string
	Cassandra  *gocql.ClusterConfig
	MongoDB    *mgo.DialInfo
}

var (
	defaultKeyspace = "default"
	defaultTarget   = "127.0.0.1"
)

func newConfig() *Config {
	return &Config{
		ID:       common.RandID("db_", 6),
		Name:     DefaultName,
		Version:  defaultVersion,
		Keyspace: defaultKeyspace,
		Target:   defaultTarget,
	}
}

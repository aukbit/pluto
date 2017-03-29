package datastore

import (
	"log"
	"regexp"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
)

// Option is used to set options for the service.
type Option interface {
	apply(*Datastore)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Datastore)

func (f optionFunc) apply(d *Datastore) {
	f(d)
}

// ID service id
func ID(id string) Option {
	return optionFunc(func(d *Datastore) {
		d.cfg.ID = id
	})
}

// Name service name
func Name(n string) Option {
	return optionFunc(func(d *Datastore) {
		d.cfg.Name = common.SafeName(n, DefaultName)
	})
}

// Keyspace db keyspace
func Keyspace(ks string) Option {
	return optionFunc(func(d *Datastore) {
		// cassandra valid characters
		//https://docs.datastax.com/en/cql/3.3/cql/cql_reference/ref-lexical-valid-chars.html
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			log.Fatal(err)
		}
		safe := reg.ReplaceAllString(ks, "_")
		d.cfg.Keyspace = safe
	})
}

// Target db address
func Target(a string) Option {
	return optionFunc(func(d *Datastore) {
		d.cfg.Target = a
	})
}

// TargetName server address
func TargetName(name string) Option {
	return optionFunc(func(d *Datastore) {
		d.cfg.TargetName = name
	})
}

// Discovery service discoery
func Discovery(dis discovery.Discovery) Option {
	return optionFunc(func(d *Datastore) {
		d.cfg.Discovery = dis
	})
}

// Logger sets a shallow copy from an input logger
func Logger(l *zap.Logger) Option {
	return optionFunc(func(d *Datastore) {
		copy := *l
		d.logger = &copy
	})
}

// Driver sets the driver to be used while establish a datastore connection
func Driver(name string) Option {
	return optionFunc(func(d *Datastore) {
		d.cfg.Driver = name
	})
}

// Cassandra sets cassandra cluster configuration
func Cassandra(cfg *gocql.ClusterConfig) Option {
	return optionFunc(func(d *Datastore) {
		d.cfg.Cassandra = cfg
		d.cfg.driver = "gocql"
	})
}

// MongoDB sets mongodb configuration
func MongoDB(cfg *mgo.DialInfo) Option {
	return optionFunc(func(d *Datastore) {
		d.cfg.MongoDB = cfg
		d.cfg.driver = "mgo"
	})
}

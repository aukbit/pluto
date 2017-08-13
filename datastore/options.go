package datastore

import (
	"github.com/aukbit/pluto/common"
	"github.com/gocql/gocql"
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
		d.cfg.Name = common.SafeName(n, defaultName)
	})
}

// Logger sets a shallow copy from an input logger
// func Logger(l *zap.Logger) Option {
// 	return optionFunc(func(d *Datastore) {
// 		copy := *l
// 		d.logger = &copy
// 	})
// }

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

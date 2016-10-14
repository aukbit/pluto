package datastore

import (
	"github.com/gocql/gocql"
	"github.com/uber-go/zap"
)

var (
	defaultCluster = "127.0.0.1"
)

type datastore struct {
	cfg     *Config
	cluster *gocql.ClusterConfig
	session *gocql.Session
	logger  zap.Logger
}

// NewServer will instantiate a new Server with the given config
func newDatastore(cfgs ...ConfigFunc) *datastore {
	c := newConfig(cfgs...)
	ds := &datastore{cfg: c, logger: zap.New(zap.NewJSONEncoder())}
	ds.setLogger()
	return ds
}

func (ds *datastore) Connect(cfgs ...ConfigFunc) {
	ds.logger.Info("connect")
	// set last configs
	for _, c := range cfgs {
		c(ds.cfg)
	}
	// set target from service discovery
	if err := ds.target(); err != nil {
		return err
	}
	ds.cluster = gocql.NewCluster(ds.cfg.Target)
	ds.cluster.ProtoVersion = 3
	ds.cluster.Keyspace = ds.cfg.Keyspace
}

func (ds *datastore) RefreshSession() error {
	ds.logger.Info("session")
	s, err := ds.cluster.CreateSession()
	if err != nil {
		return err
	}
	ds.session = s
	return nil
}

func (ds *datastore) Close() {
	ds.logger.Info("close")
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

func (ds *datastore) setLogger() {
	ds.logger = ds.logger.With(
		zap.Nest("cassandra",
			zap.String("target", ds.cfg.Target),
			zap.String("keyspace", ds.cfg.Keyspace)))
}

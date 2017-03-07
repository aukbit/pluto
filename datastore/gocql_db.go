package datastore

import (
	"context"

	"go.uber.org/zap"

	"github.com/gocql/gocql"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	defaultCluster = "127.0.0.1"
)

type datastore struct {
	cfg     *Config
	cluster *gocql.ClusterConfig
	session *gocql.Session
	logger  *zap.Logger
	health  *health.Server
}

// NewServer will instantiate a new Server with the given config
func newDatastore(cfgs ...ConfigFunc) *datastore {
	c := newConfig(cfgs...)
	d := &datastore{cfg: c,
		health: health.NewServer(),
	}
	d.logger, _ = zap.NewProduction()
	return d
}

func (ds *datastore) Connect(cfgs ...ConfigFunc) error {
	// set last configs
	for _, c := range cfgs {
		c(ds.cfg)
	}
	// register at service discovery
	if err := ds.register(); err != nil {
		return err
	}
	// set target from service discovery
	if err := ds.target(); err != nil {
		return err
	}
	// set logger
	ds.setLogger()
	ds.logger.Info("connect")
	ds.cluster = gocql.NewCluster(ds.cfg.Target)
	ds.cluster.ProtoVersion = 3
	ds.cluster.Keyspace = ds.cfg.Keyspace
	// set health
	ds.health.SetServingStatus(ds.cfg.ID, 1)
	return nil
}

func (ds *datastore) RefreshSession() error {
	ds.logger.Info("session")
	s, err := ds.cluster.CreateSession()
	if err != nil {
		ds.health.SetServingStatus(ds.cfg.ID, 2)
		return err
	}
	ds.session = s
	ds.health.SetServingStatus(ds.cfg.ID, 1)
	return nil
}

func (ds *datastore) Config() *Config {
	return ds.cfg
}

func (ds *datastore) Close() {
	ds.logger.Info("close")
	// set health as not serving
	ds.health.SetServingStatus(ds.cfg.ID, 2)
	// unregister
	ds.unregister()
	ds.session.Close()
}

func (ds *datastore) Session() *gocql.Session {
	return ds.session
}

func (ds *datastore) Health() *healthpb.HealthCheckResponse {
	ds.RefreshSession()
	hcr, err := ds.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: ds.cfg.ID})
	if err != nil {
		ds.logger.Error("Health", zap.String("err", err.Error()))
		return &healthpb.HealthCheckResponse{Status: 2}
	}
	return hcr
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
		zap.String("type", "db"),
		zap.String("id", ds.cfg.ID),
		zap.String("name", ds.cfg.Name),
		zap.String("target", ds.cfg.Target),
		zap.String("keyspace", ds.cfg.Keyspace))
}

package datastore

import (
	"context"
	"fmt"

	mgo "gopkg.in/mgo.v2"

	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// DefaultName prefix datastore client name
	DefaultName    = "client_db"
	defaultVersion = "v1.0.0"
)

var (
	defaultCluster = "127.0.0.1"
)

type Datastore struct {
	cfg     *Config
	session interface{}
	logger  *zap.Logger
	health  *health.Server
}

// New creates a default datatore
func New(opts ...Option) *Datastore {
	return newDatastore(opts...)
}

// NewServer will instantiate a new Server with the given config
func newDatastore(opts ...Option) *Datastore {
	d := &Datastore{
		cfg:    newConfig(),
		health: health.NewServer(),
	}
	// d.logger, _ = zap.NewProduction()
	if len(opts) > 0 {
		d = d.WithOptions(opts...)
	}
	return d
}

// WithOptions clones the current Client, applies the supplied Options, and
// returns the resulting Client. It's safe to use concurrently.
func (ds *Datastore) WithOptions(opts ...Option) *Datastore {
	c := ds.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// clone creates a shallow copy client
func (ds *Datastore) clone() *Datastore {
	copy := *ds
	return &copy
}

func (ds *Datastore) Connect(opts ...Option) error {
	// set last configs
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(ds)
		}
	}
	// set logger
	ds.setLogger()
	// create db session
	switch ds.cfg.driver {
	case "gocql":
		s, err := gocql.NewSession(cfg)
		if err != nil {
			return err
		}
	case "mgo":
		s, err := mgo.DialWithInfo(ds.cfg.MongoDB)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("datastore driver not available %v", ds.cfg.Driver)
	}
	ds.session = s
	ds.logger.Info("connected")
	// set health
	ds.health.SetServingStatus(ds.cfg.ID, healthpb.HealthCheckResponse_SERVING)
	return nil
}

func (ds *Datastore) RefreshSession() error {
	ds.logger.Info("session")
	switch ds.cfg.Driver {
	case "gocql":
		f := gocql.NewCluster(hosts)
		ds.session.(*gocql.Session)
		s, err := ds.cluster.CreateSession()
		if err != nil {
			ds.health.SetServingStatus(ds.cfg.ID, 2)
			return err
		}
	}
	ds.session = s
	ds.health.SetServingStatus(ds.cfg.ID, healthpb.HealthCheckResponse_SERVING)
	return nil
}

func (ds *Datastore) Close() {
	ds.logger.Info("close")
	// set health as not serving
	ds.health.SetServingStatus(ds.cfg.ID, healthpb.HealthCheckResponse_NOT_SERVING)
	//
	switch ds.cfg.Driver {
	case "gocql":
		ds.session.(*gocql.Session).Close()
	case "mgo":
		ds.session.(*mgo.Session).Close()
	default:
		//
	}
}

func (ds *Datastore) Session() interface{} {
	switch ds.cfg.Driver {
	case "gocql":
		ds.NewCluster(ds.cfg.Keyspace, ds.cfg.Target)
	default:
		return fmt.Errorf("datastore driver not available %v", ds.cfg.Driver)
	}
	return ds.session
}

func (ds *Datastore) Health() *healthpb.HealthCheckResponse {
	ds.RefreshSession()
	hcr, err := ds.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: ds.cfg.ID})
	if err != nil {
		ds.logger.Error("Health", zap.String("err", err.Error()))
		return &healthpb.HealthCheckResponse{Status: 2}
	}
	return hcr
}

func (ds *Datastore) createKeyspace(keyspace string, replicationFactor int) error {
	q := "CREATE KEYSPACE ? WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : ? };"
	if err := ds.session.Query(q, keyspace, replicationFactor).Exec(); err != nil {
		return err
	}
	return nil
}

func (ds *Datastore) setLogger() {
	ds.logger = ds.logger.With(
		zap.String("type", "db"),
		zap.String("id", ds.cfg.ID),
		zap.String("name", ds.cfg.Name),
		zap.String("target", ds.cfg.Target),
		zap.String("keyspace", ds.cfg.Keyspace))
}

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
	defaultName    = "db"
	defaultVersion = "1.3.0"
)

type Datastore struct {
	cfg    *Config
	logger *zap.Logger
	health *health.Server
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

func (ds *Datastore) Init(opts ...Option) error {
	// set last configs
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(ds)
		}
	}
	// set logger
	ds.setLogger()
	ds.logger.Info("init")
	s, err := ds.NewSession()
	if err != nil {
		return err
	}
	defer ds.Close(s)
	// set health
	ds.health.SetServingStatus(ds.cfg.ID, healthpb.HealthCheckResponse_SERVING)
	return nil
}

func (ds *Datastore) Close(session interface{}) {
	switch session.(type) {
	case *gocql.Session:
		defer session.(*gocql.Session).Close()
	case *mgo.Session:
		defer session.(*mgo.Session).Close()
	}
}

func (ds *Datastore) NewSession() (session interface{}, err error) {
	switch ds.cfg.driver {
	case "gocql":
		session, err = gocql.NewSession(*ds.cfg.Cassandra)
		if err != nil {
			return nil, err
		}
	case "mgo":
		session, err = mgo.DialWithInfo(ds.cfg.MongoDB)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("datastore driver not available")
	}
	return session, nil
}

func (ds *Datastore) setHealth() {
	session, err := ds.NewSession()
	if err != nil {
		ds.health.SetServingStatus(ds.cfg.ID, healthpb.HealthCheckResponse_NOT_SERVING)
	}
	defer ds.Close(session)
	ds.health.SetServingStatus(ds.cfg.ID, healthpb.HealthCheckResponse_SERVING)
}

func (ds *Datastore) Health() *healthpb.HealthCheckResponse {
	ds.setHealth()
	hcr, err := ds.health.Check(
		context.Background(),
		&healthpb.HealthCheckRequest{Service: ds.cfg.ID},
	)
	if err != nil {
		ds.logger.Error("Health", zap.String("err", err.Error()))
		return &healthpb.HealthCheckResponse{
			Status: healthpb.HealthCheckResponse_NOT_SERVING,
		}
	}
	return hcr
}

func (ds *Datastore) setLogger() {
	ds.logger = ds.logger.With(
		zap.String("id", ds.cfg.ID),
		zap.String("name", ds.cfg.Name),
		zap.String("driver", ds.cfg.driver),
	)
}

// Name returns datastore name
func (ds *Datastore) Name() string {
	return ds.cfg.Name
}

// IsAvailable checks if a Datastore instance as been initialized
func (ds *Datastore) IsAvailable() bool {
	return ds != nil
}

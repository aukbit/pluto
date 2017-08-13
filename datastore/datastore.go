package datastore

import (
	"context"
	"fmt"
	"os"

	mgo "gopkg.in/mgo.v2"

	"github.com/gocql/gocql"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	// DefaultName prefix datastore client name
	defaultName = "db"
)

// contextKey datstore context keys
type contextKey string

type Datastore struct {
	cfg    Config
	health *health.Server
	logger zerolog.Logger
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
	d.logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
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
	ds.logger.With().Str("id", ds.cfg.ID).Str("name", ds.cfg.Name).Str("driver", ds.cfg.driver)
	// set last configs
	if len(opts) > 0 {
		for _, opt := range opts {
			opt.apply(ds)
		}
	}
	s, err := ds.NewSession()
	if err != nil {
		return err
	}
	defer ds.Close(s)
	// set health
	ds.health.SetServingStatus(ds.cfg.ID, healthpb.HealthCheckResponse_SERVING)
	ds.logger.Info().Msg(fmt.Sprintf("%s initialized", ds.Name()))
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

// WithContext add session to existing context
func WithContext(ctx context.Context, session interface{}) context.Context {
	return context.WithValue(ctx, contextKey("session"), session)
}

// FromContext returns datastore session instance from a context
func FromContext(ctx context.Context) interface{} {
	return ctx.Value(contextKey("session"))
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
		ds.logger.Error().Msg(fmt.Sprintf("%s Health() %v", ds.Name(), err.Error()))
		return &healthpb.HealthCheckResponse{
			Status: healthpb.HealthCheckResponse_NOT_SERVING,
		}
	}
	return hcr
}

// Name returns datastore name
func (ds *Datastore) Name() string {
	return ds.cfg.Name
}

// IsAvailable checks if a Datastore instance as been initialized
func (ds *Datastore) IsAvailable() bool {
	return ds != nil
}

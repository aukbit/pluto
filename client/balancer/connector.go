package balancer

import (
	"context"

	g "github.com/aukbit/pluto/client/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// type Connector interface {
// 	Client() interface{}
// 	Connector() *Connector
// 	Health() bool
// }

const (
	// DefaultName prefix connector name
	defaultName    = "connector"
	defaultVersion = "v1.0.0"
)

type ConnsCh chan *Connector

// Connector struct
type Connector struct {
	cfg        *Config
	requestsCh chan Request          // requests channel to receive requests from balancer
	pending    int                   // count pending tasks
	index      int                   // index in the heap
	conn       *grpc.ClientConn      // grpc connection to communicate with the server
	client     interface{}           // grpc client stub to perform RPCs
	stopCh     chan bool             // receive a stop call
	doneCh     chan bool             // guarantees has beeen stopped correctly
	health     healthpb.HealthClient // Client API for Health service
	logger     *zap.Logger
}

// NewConnector returns a new connector with options passed in
func NewConnector(opts ...Option) *Connector {
	return newConnector(opts...)
}

// newConnector ...
func newConnector(opts ...Option) *Connector {
	c := &Connector{
		cfg:        newConfig(),
		requestsCh: make(chan Request),
		stopCh:     make(chan bool),
		doneCh:     make(chan bool),
	}
	c.logger, _ = zap.NewProduction()
	if len(opts) > 0 {
		c = c.WithOptions(opts...)
	}
	return c
}

// WithOptions clones the current Client, applies the supplied Options, and
// returns the resulting Client. It's safe to use concurrently.
func (c *Connector) WithOptions(opts ...Option) *Connector {
	d := c.clone()
	for _, opt := range opts {
		opt.apply(d)
	}
	return d
}

// clone creates a shallow copy client
func (c *Connector) clone() *Connector {
	copy := *c
	return &copy
}

// dial establish grpc client connection with the grpc server
func (c *Connector) dial() error {
	c.logger.Info("dial")
	// append logger
	c.cfg.UnaryClientInterceptors = append(c.cfg.UnaryClientInterceptors, loggerUnaryClientInterceptor(c))
	// dial
	// TODO use TLS
	conn, err := grpc.Dial(
		c.cfg.Target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(g.WrapperUnaryClient(c.cfg.UnaryClientInterceptors...)))
	if err != nil {
		return err
	}
	// keep connection for later close
	c.conn = conn
	// register proto client to get a stub to perform RPCs
	c.client = c.cfg.GRPCRegister(conn)
	// register proto health client to get a stub to perform RPCs
	c.health = healthpb.NewHealthClient(conn)
	return nil
}

// watch waits for any call from balancer
func (c *Connector) watch() {
	c.logger.Info("watch")
	for {
		select {
		case req := <-c.requestsCh: // get request from balancer
			req.connsCh <- c
		case <-c.stopCh:
			close(c.doneCh)
			return
		}
	}
}

func (c *Connector) Init() error {
	c.logger = c.logger.With(
		zap.String("target", c.cfg.Target),
	)
	if err := c.dial(); err != nil {
		return err
	}
	go c.watch()
	return nil
}

func (c *Connector) Client() interface{} {
	return c.client
}

func (c *Connector) Connector() *Connector {
	return c
}

// Close stops connector and close grpc connection
func (c *Connector) Close() {
	c.conn.Close()
	c.stopCh <- true
	<-c.doneCh
	c.logger.Info("closed")
}

// Health check if a round trip with server is valid or not
func (c *Connector) Health() bool {
	hcr, err := c.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		return false
	}
	return hcr.Status.String() == "SERVING"
}

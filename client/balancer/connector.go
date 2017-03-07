package balancer

import (
	"context"

	g "github.com/aukbit/pluto/client/grpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Connector interface {
	Client() interface{}
	Connector() *connector
	Health() bool
}

const (
	// DefaultName prefix connector name
	DefaultName    = "connector"
	defaultVersion = "v1.0.0"
)

// NewConnector returns a new connector with cfg passed in
func NewConnector(cfgs ...ConfigFn) *connector {
	return newConnector(cfgs...)
}

type ConnsCh chan *connector

// connector struct
type connector struct {
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

// newConnector ...
func newConnector(cfgs ...ConfigFn) *connector {
	c := newConfig(cfgs...)
	conn := &connector{
		cfg:        c,
		requestsCh: make(chan Request),
		stopCh:     make(chan bool),
		doneCh:     make(chan bool),
	}
	conn.logger, _ = zap.NewProduction()
	conn.initLogger()
	return conn
}

// dial establish grpc client connection with the grpc server
func (c *connector) dial() error {
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
func (c *connector) watch() {
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

func (c *connector) Init() error {
	if err := c.dial(); err != nil {
		return err
	}
	go c.watch()
	return nil
}

func (c *connector) Client() interface{} {
	return c.client
}

func (c *connector) Connector() *connector {
	return c
}

// Close stops connector and close grpc connection
func (c *connector) Close() {
	c.logger.Info("close")
	c.conn.Close()
	c.stopCh <- true
	<-c.doneCh
}

// Health check if a round trip with server is valid or not
func (c *connector) Health() bool {
	hcr, err := c.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{})
	if err != nil {
		return false
	}
	return hcr.Status.String() == "SERVING"
}

func (c *connector) initLogger() {
	c.logger = c.logger.With(
		zap.String("type", "connector"),
		zap.String("id", c.cfg.ID),
		zap.String("name", c.cfg.Name),
		zap.String("target", c.cfg.Target),
		zap.String("parent", c.cfg.ParentID),
	)
}

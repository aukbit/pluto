package balancer

import (
	"context"

	"github.com/uber-go/zap"

	g "bitbucket.org/aukbit/pluto/client/grpc"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Connector channel of requests
type Connector struct {
	cfg        *Config
	requestsCh chan Request          // requests channel to receive requests from balancer
	pending    int                   // count pending tasks
	index      int                   // index in the heap
	conn       *grpc.ClientConn      // grpc connection to communicate with the server
	Client     interface{}           // grpc client stub to perform RPCs
	stopCh     chan bool             // receive a stop call
	doneCh     chan bool             // guarantees has beeen stopped correctly
	health     healthpb.HealthClient // Client API for Health service
	logger     zap.Logger
}

// NewConnector ...
func NewConnector(cfgs ...ConfigFn) *Connector {
	c := newConfig(cfgs...)
	conn := &Connector{
		cfg:        c,
		requestsCh: make(chan Request),
		stopCh:     make(chan bool),
		doneCh:     make(chan bool),
		logger:     zap.New(zap.NewJSONEncoder())}

	conn.initLogger()
	return conn
}

// Dial establish grpc client connection with the grpc server
func (c *Connector) Dial() error {
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
	c.Client = c.cfg.GRPCRegister(conn)
	// register proto health client to get a stub to perform RPCs
	c.health = healthpb.NewHealthClient(conn)
	return nil
}

// Watch waits for any call from balancer
func (c *Connector) Watch() {
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

// Close stops connector and close grpc connection
func (c *Connector) Close() {
	c.logger.Info("close")
	c.conn.Close()
	c.stopCh <- true
	<-c.doneCh
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

func (c *Connector) initLogger() {
	c.logger = c.logger.With(
		zap.Nest("connector",
			zap.String("id", c.cfg.ID),
			zap.String("name", c.cfg.Name),
			zap.String("target", c.cfg.Target),
			zap.String("parent", c.cfg.ParentID)))
}

package balancer

import (
	"github.com/uber-go/zap"

	g "bitbucket.org/aukbit/pluto/client/grpc"
	"google.golang.org/grpc"
)

// Connector channel of requests
type Connector struct {
	cfg        *Config
	requestsCh chan Request     // requests channel to receive requests from balancer
	pending    int              // count pending tasks
	index      int              // index in the heap
	conn       *grpc.ClientConn // grpc connection to communicate with the server
	Client     interface{}      // grpc client stub to perform RPCs
	stopCh     chan bool        // receive a stop call
	doneCh     chan bool        // guarantees has beeen stopped correctly
	logger     zap.Logger
}

// NewConnector ...
func NewConnector(cfgs ...ConfigFn) *Connector {
	c := newConfig(cfgs...)
	return &Connector{
		cfg:        c,
		requestsCh: make(chan Request),
		stopCh:     make(chan bool),
		doneCh:     make(chan bool),
		logger:     zap.New(zap.NewJSONEncoder())}
}

// Dial establish grpc client connection with the grpc server
func (c *Connector) Dial() error {
	c.logger.Info("dial")
	// TODO use TLS
	// append logger
	// c.cfg.UnaryClientInterceptors = append(c.cfg.UnaryClientInterceptors, loggerUnaryClientInterceptor(c))
	// dial
	conn, err := grpc.Dial(
		c.cfg.Target,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(g.WrapperUnaryClient(c.cfg.UnaryClientInterceptors...)))
	if err != nil {
		return err
	}
	// keep connection for later close
	c.conn = conn
	// register health methods on connection
	// c.healthCall = healthpb.NewHealthClient(conn)
	// register proto client to get a stub to perform RPCs
	c.Client = c.cfg.GRPCRegister(conn)
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

// Stop stops connector and close grpc connection
func (c *Connector) Stop() {
	c.conn.Close()
	c.stopCh <- true
	<-c.doneCh
}

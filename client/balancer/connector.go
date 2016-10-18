package balancer

import (
	"github.com/uber-go/zap"

	"google.golang.org/grpc"
)

// Connector channel of requests
type Connector struct {
	cfg        *Config
	requestsCh chan Request     // requests channel to receive requests from balancer
	pending    int              // count pending tasks
	index      int              // index in the heap
	conn       *grpc.ClientConn // grpc connection to communicate with the server
	client     interface{}      // grpc client stub to perform RPCs
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

// dial establish client grpc connection with the grpc server
func (c *Connector) dial() error {
	conn, err := grpc.Dial(
		c.cfg.Target,
		grpc.WithInsecure())

	if err != nil {
		c.logger.Error("dial", zap.String("err", err.Error()))
		return err
	}
	// keep connection for later close
	c.conn = conn
	// register proto client to get a stub to perform RPCs
	c.client = c.cfg.GRPCRegister(conn)
	return nil
}

// watch waits for any call from balancer
func (c *Connector) watch() {
	for {
		select {
		case req := <-c.requestsCh: // get request from balancer
			req.connsCh <- c
		case <-c.stopCh:
			c.conn.Close()
			close(c.doneCh)
			return
		}
	}
}

func (c *Connector) stop() {
	c.stopCh <- true
	<-c.doneCh
}

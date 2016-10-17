package balancer

import (
	"log"

	"google.golang.org/grpc"
)

// Call the caller sends calls to the balancer
type Call struct {
	clientsCh chan interface{} // grpc client stub channel
	conn      *Connector
}

// done receives a callsCh channel with write only constraint
func (c *Call) done(doneCh chan<- *Connector) {
	log.Printf("Call it's done") //
	doneCh <- c.conn
}

// Connector channel of requests
type Connector struct {
	callsCh chan *Call       // call channel to receive calls from balancer
	pending int              // count pending tasks
	index   int              // index in the heap
	target  string           // grpc server address
	conn    *grpc.ClientConn // grpc connection to communicate with the server
	client  interface{}      // grpc client stub to perform RPCs
}

// dial establish client grpc connection with the grpc server
func (c *Connector) dial(regFn func(*grpc.ClientConn) interface{}) error {
	log.Printf("Connector dial")
	conn, err := grpc.Dial(
		c.target,
		grpc.WithInsecure())

	if err != nil {
		// dc.logger.Error("dial", zap.String("err", err.Error()))
		log.Fatalf("dial %v", err)
		return err
	}
	// keep connection for later close
	c.conn = conn
	// register proto client to get a stub to perform RPCs
	c.client = regFn(conn)
	return nil
}

// watch waits for any call from balancer
func (c *Connector) watch() {
	log.Printf("Connector watch")
	// doneCh := make(chan *Connector)
	for {
		select {
		case call := <-c.callsCh: // get call from balancer
			log.Printf("call received")
			call.clientsCh <- c.client // send client stub over the call clients channel
			// set connector reference in the call
			call.conn = c
			// case <-c.outCh: // call it's done
			// 	c.done(doneCh)
		}
	}
}

// done receives a connectors channel with write only constraint
func (c *Connector) done(doneCh chan<- *Connector) {
	log.Printf("Connector it's done")
	doneCh <- c // send the connector over the channel to inform balancer the request is finish
}

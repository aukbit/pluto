package balancer

import (
	"log"

	"google.golang.org/grpc"
)

// Call the caller sends calls to the balancer
type Call struct {
	clientsCh chan interface{} // grpc client stub channel
}

// func (c *Call) watch(conn *Connector, doneCh chan<- *Connector) {
// 	log.Printf("Call watch")
// 	<-c.doneCh
// 	doneCh <- conn
// }
//
// // done
// func (c *Call) done() {
// 	log.Printf("Call it's done") //
// 	c.doneCh <- true
// }

// Connector channel of requests
type Connector struct {
	callsCh chan *Call       // call channel to receive calls from balancer
	pending int              // count pending tasks
	index   int              // index in the heap
	target  string           // grpc server address
	conn    *grpc.ClientConn // grpc connection to communicate with the server
	client  interface{}      // grpc client stub to perform RPCs
	doneCh  chan bool        //
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
func (c *Connector) watch(doneCh chan<- *Connector) {
	log.Printf("Connector watch")
	for {
		select {
		case callCh := <-c.callsCh: // get call from balancer
			log.Printf("call received %v", c.target) //
			call.clientsCh <- c.client               // send client stub over the call clients channel
			// wait for call to finish
			// go call.watch(c, doneCh)

		case <-c.doneCh:
			doneCh <- c
		}
	}
}

// done receives a connectors channel with write only constraint
func (c *Connector) done() {
	log.Printf("Connector it's done")
	c.doneCh <- true
}

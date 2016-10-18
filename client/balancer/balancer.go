package balancer

import (
	"container/heap"
	"log"
)

// Balancer needs a pool of connectors and a single channel to which requesters
// can report task completion
type Balancer struct {
	pool   Pool
	doneCh chan *Connector
}

// balance receives calls on call channel with read only constraint
func (b *Balancer) balance(inCh <-chan *Call) {
	log.Printf("balance")
	for {
		select {
		case call := <-inCh: // received a callCh
			log.Printf("balance: received a callCh: %v", call)
			b.dispatch(call) //send it to a Connector
		case c := <-b.doneCh: // a call has finished
			log.Printf("balance: a call has finished %v", c.index)
			b.completed(c)
		}
	}
}

func (b *Balancer) dispatch(call *Call) {
	log.Printf("dispatch call: %v", call)
	// get the least loaded connector..
	c := heap.Pop(&b.pool).(*Connector)
	// send it the call
	c.callsCh <- call
	// one more in its work queue
	c.pending++
	// put it into its place on the heap
	heap.Push(&b.pool, c)
}

func (b *Balancer) completed(c *Connector) {
	log.Printf("completed %v", c)
	// remove one from the queue
	c.pending--
	// remove it from the heap
	heap.Remove(&b.pool, c.index)
	// put it into its place on the heap
	heap.Push(&b.pool, c)
}

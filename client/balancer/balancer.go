package balancer

import "container/heap"

// Balancer needs a pool of connectors and a single channel to which requesters
// can report task completion
type Balancer struct {
	pool    Pool
	connsCh chan *Connector
}

// NewBalancer starts a balancer with an empty pool
func NewBalancer() *Balancer {
	return &Balancer{
		pool:    newPool(),
		connsCh: make(chan *Connector)}
}

// balance receives requests with read only constraint
func (b *Balancer) balance(inCh <-chan Request) {
	for {
		select {
		case req := <-inCh: // received a request
			b.dispatch(req) //send it to a Connector
		case c := <-b.connsCh: // a request has finished
			b.completed(c) //
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	// get the least loaded connector..
	c := heap.Pop(&b.pool).(*Connector)
	// send it the call
	c.requestsCh <- req
	// one more in its work queue
	c.pending++
	// put it into its place on the heap
	heap.Push(&b.pool, c)
}

func (b *Balancer) completed(c *Connector) {
	// remove one from the queue
	c.pending--
	// remove it from the heap
	heap.Remove(&b.pool, c.index)
	// put it into its place on the heap
	heap.Push(&b.pool, c)
}

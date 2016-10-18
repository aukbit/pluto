package balancer

import "container/heap"

// Balancer needs a pool of connectors and a single channel to which requesters
// can report task completion
type Balancer struct {
	pool    Pool
	ConnsCh chan *Connector
}

// NewBalancer starts a balancer with an empty pool
func NewBalancer() *Balancer {
	return &Balancer{
		pool:    newPool(),
		ConnsCh: make(chan *Connector)}
}

// Push pushes the connector onto the heap
func (b *Balancer) Push(c *Connector) {
	heap.Push(&b.pool, c)
}

// Pop removes the minimum element (according to Less)
// from the heap and returns it
func (b *Balancer) Pop() *Connector {
	return heap.Pop(&b.pool).(*Connector)
}

func (b *Balancer) Pool() Pool {
	return b.pool
}

// balance receives requests with read only constraint
func (b *Balancer) Balance(ch <-chan Request) {
	for {
		select {
		case req := <-ch: // received a request
			b.dispatch(req) //send it to a Connector
		case c := <-b.ConnsCh: // a request has finished
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

package balancer

import "container/heap"

// Balancer needs a pool of connectors and a single channel to which requesters
// can report task completion
type Balancer struct {
	pool    Pool
	connsCh ConnsCh
}

// NewBalancer starts a balancer with an empty pool
func NewBalancer() *Balancer {
	return &Balancer{
		pool:    newPool(),
		connsCh: make(chan *connector)}
}

// Push pushes the connector onto the heap
func (b *Balancer) Push(c *connector) {
	heap.Push(&b.pool, c)
}

// Pop removes the minimum element (according to Less)
// from the heap and returns it
func (b *Balancer) Pop() *connector {
	return heap.Pop(&b.pool).(*connector)
}

// Pool returns balancer pool
func (b *Balancer) Pool() Pool {
	return b.pool
}

// Done send connector over connsCh channel
func (b *Balancer) Done(conn *connector) {
	// send conn over balancer connsCh
	b.connsCh <- conn
}

// Balance receives requests with read only constraint
func (b *Balancer) Balance(ch <-chan Request) {
	for {
		select {
		case req := <-ch: // received a request
			b.dispatch(req) //send it to aconnector
		case c := <-b.connsCh: // a request has finished
			b.completed(c) //
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	// get the least loaded connector..
	c := heap.Pop(&b.pool).(*connector)
	// send it the call
	c.requestsCh <- req
	// one more in its work queue
	c.pending++
	// put it into its place on the heap
	heap.Push(&b.pool, c)
}

func (b *Balancer) completed(c *connector) {
	// remove one from the queue
	c.pending--
	// remove it from the heap
	heap.Remove(&b.pool, c.index)
	// put it into its place on the heap
	heap.Push(&b.pool, c)
}

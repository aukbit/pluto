package balancer

import (
	"container/heap"
	"log"
)

// Request The requester sends requests to the balancer
type Request struct {
	fn func() int // The operation to perform
	c  chan int   // The channel to return the result
}

// Balancer needs a pool of workers and a single channel to which requesters
// can report task completion
type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) balance(work chan Request) {
	log.Printf("balance workCh: %v", work)
	for {
		select {
		case req := <-work: // received a Request
			log.Printf("balance: Received a Request req: %v", req)
			b.dispatch(req) //send it to a Worker
		case w := <-b.done: // a worker has finished
			log.Printf("balance: a worker has finished %v", w.index)
			b.completed(w)
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	log.Printf("dispatch 1 req: %v", req)
	// get the least loaded worker..
	w := heap.Pop(&b.pool).(*Worker)
	log.Printf("dispatch 2 worker: %v", w)
	// send it the task
	w.requests <- req
	// one more in its work queue
	w.pending++
	log.Printf("dispatch 3 %v", w.pending)
	// put it into its place on the heap
	heap.Push(&b.pool, w)
}

func (b *Balancer) completed(w *Worker) {
	log.Printf("completed %v", w)
	// remove one from the queue
	w.pending--
	// remove it from the heap
	heap.Remove(&b.pool, w.index)
	// put it into its place on the heap
	heap.Push(&b.pool, w)
}

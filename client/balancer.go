package balancer

import (
	"container/heap"
	"log"
)

// Balancer needs a pool of workers and a single channel to which requesters
// can report task completion
type Balancer struct {
	pool   Pool
	doneCh chan *Worker
}

func (b *Balancer) balance(workCh chan Request) {
	log.Printf("balance workCh: %v", workCh)
	for {
		select {
		case req := <-workCh: // received a Request
			log.Printf("balance: received a Request req: %v", req)
			b.dispatch(req) //send it to a Worker
		case w := <-b.doneCh: // a worker has finished
			log.Printf("balance: a worker has finished %v", w.index)
			b.completed(w)
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	log.Printf("dispatch req: %v", req)
	// get the least loaded worker..
	w := heap.Pop(&b.pool).(*Worker)
	// send it the task
	w.requestsCh <- req
	// one more in its work queue
	w.pending++
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

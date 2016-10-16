package balancer

import "log"

// Worker channel of requests
type Worker struct {
	requests chan Request // work to do
	pending  int          // count pending tasks
	index    int          // index in the heap
}

func (w *Worker) work(done chan *Worker) {
	log.Printf("Worker work %v", done)
	for {
		req := <-w.requests // get Request from balancer
		req.c <- req.fn()   // call fn and send result
		done <- w           // this request is finish
	}
}

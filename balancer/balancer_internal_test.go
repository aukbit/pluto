package balancer

import (
	"log"
	"testing"
	"time"
)

var NumWorkers = 1

func workFn() int {
	log.Printf("Hello workFn")
	time.Sleep(1 * time.Second)
	return 121
}

func RunBalancer(work chan Request) {
	wA := &Worker{requests: make(chan Request), pending: 0}
	log.Printf("RunBalancer worker: %v", wA)
	// wB := &Worker{requests: make(chan Request), pending: 0}
	p := Pool{wA}
	log.Printf("RunBalancer pool: %v", p)
	b := &Balancer{pool: p, done: make(chan *Worker)}
	log.Printf("RunBalancer balancer: %v", b)
	b.balance(work)
}

func requester(work chan<- Request) {
	log.Printf("requester workCh: %v", work)
	c := make(chan int)
	for {
		// Kill some time (fake load).
		log.Printf("requester sleep")
		time.Sleep(time.Second * 3)
		log.Printf("requester sleep over")
		req := Request{workFn, c}
		work <- req // send request
		log.Printf("requester send request workCh:%v req: %v", work, req)
		result := <-c // wait for answer
		// furtherProcess(result)
		log.Printf("requester result %v", result)
	}
}

func TestBalancer(t *testing.T) {
	workCh := make(chan Request)
	log.Printf("TestBalancer reqChan: %v", workCh)
	for i := 0; i < NumWorkers; i++ {
		go requester(workCh)
	}
	RunBalancer(workCh)
	// receiveLotsOfResults(out)
	//
	// go RunBalancer()
	//
	// r := make(chan Request)
	// requester(r)

}

package balancer

import (
	"log"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup
var NumWorkers = 2

func workFn() int {
	log.Printf("Hello workFn")
	time.Sleep(1 * time.Second)
	return 121
}

func RunBalancer(work chan Request) {
	defer wg.Done()
	doneCh := make(chan *Worker)
	requestsCh := make(chan Request)
	wA := &Worker{requests: requestsCh, pending: 0}
	wB := &Worker{requests: requestsCh, pending: 0}
	go wA.work(doneCh)
	go wB.work(doneCh)
	p := Pool{wA, wB}
	log.Printf("RunBalancer pool: %v", p)
	b := &Balancer{pool: p, done: doneCh}
	log.Printf("RunBalancer balancer: %v", b)
	b.balance(work)
}

func requester(work chan<- Request) {
	defer wg.Done()
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
		wg.Add(1)
		go requester(workCh)
	}
	wg.Add(1)
	go RunBalancer(workCh)
	wg.Wait()
	log.Printf("TestBalancer END")
	// receiveLotsOfResults(out)
	//
	// go RunBalancer()
	//
	// r := make(chan Request)
	// requester(r)

}

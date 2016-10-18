package balancer

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"bitbucket.org/aukbit/pluto/server"
	pb "bitbucket.org/aukbit/pluto/server/proto"
	"github.com/paulormart/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v", in.Name)}, nil
}

func TestMain(m *testing.M) {
	// Create pluto server
	s1 := server.NewServer(
		server.Name("test_gopher_1"),
		server.Addr(":65070"),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}))
	s2 := server.NewServer(
		server.Name("test_gopher_2"),
		server.Addr(":65071"),
		server.GRPCRegister(func(g *grpc.Server) {
			pb.RegisterGreeterServer(g, &greeter{})
		}))

	if !testing.Short() {
		// Run Server
		go func() {
			if err := s1.Run(); err != nil {
				log.Fatal(err)
			}
		}()
		time.Sleep(time.Millisecond * 100)
		go func() {
			if err := s2.Run(); err != nil {
				log.Fatal(err)
			}
		}()
		time.Sleep(time.Millisecond * 100)
	}
	result := m.Run()
	if !testing.Short() {
		// Stop Server
		s1.Stop()
		time.Sleep(time.Millisecond * 100)
		s2.Stop()
		time.Sleep(time.Millisecond * 100)
	}
	os.Exit(result)
}

func InitConnectors(callsCh chan *Call) (cons []*Connector) {
	cA := &Connector{target: "localhost:65070", callsCh: callsCh}
	cB := &Connector{target: "localhost:65071", callsCh: callsCh}
	// establish connectors
	cA.dial(func(cc *grpc.ClientConn) interface{} {
		return pb.NewGreeterClient(cc)
	})
	cB.dial(func(cc *grpc.ClientConn) interface{} {
		return pb.NewGreeterClient(cc)
	})
	// watch for requests
	go cA.watch(doneCh)
	go cB.watch(doneCh)

	cons = append(cons, cA, cB)
	return cons
}

func InitBalancer(cons []*Connector, inCh <-chan *Call, outCh chan<- *Call, doneCh chan *Connector) {
	// initialize pool with connectors available
	p := Pool{}
	for _, c := range cons {
		p.Push(c)
	}
	heap.Init(&p)
	log.Printf("RunBalancer pool: %v", p)
	// set balancer

	b := &Balancer{pool: p, doneCh: doneCh}
	log.Printf("RunBalancer balancer: %v", b)
	b.balance(inCh)
}

func TestBalancer(t *testing.T) {
	callsCh := make(chan *Call)
	cons := InitConnectors(callsCh)
	go InitBalancer(cons, callsCh)

	// connectors channel
	connsCh := make(chan *Connector)
	c := &Call{connsCh: connsCh}
	time.Sleep(time.Second * 1)
	log.Printf("send call over channel")
	callsCh <- c

	// wait for client
	client := <-clientsCh
	// Make a Call
	r, err := client.(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: "Gopher"})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "Hello Gopher", r.Message)
	// set call has done
	c.done()

	time.Sleep(time.Second * 1)
	log.Printf("TestBalancer END")
}

// func requester(work chan<- Request) {
// 	log.Printf("requester workCh: %v", work)
// 	c := make(chan int)
// 	for {
// 		// Kill some time (fake load).
// 		time.Sleep(time.Second * 3)
// 		req := Request{workFn, c}
// 		work <- req // send request
// 		log.Printf("requester send request workCh:%v req: %v", work, req)
// 		result := <-c // wait for answer
// 		// furtherProcess(result)
// 		log.Printf("requester result %v", result)
// 	}
// }

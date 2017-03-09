package balancer

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/paulormart/assert"

	"github.com/aukbit/pluto/server"
	pb "github.com/aukbit/pluto/test/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var numServers = 1
var numRequests = 1

const PORT = 65070

var wg sync.WaitGroup

type greeter struct{}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	t := rand.Int31n(1000)
	time.Sleep(time.Duration(t) * time.Millisecond)
	return &pb.HelloReply{Message: fmt.Sprintf("Hello %v", in.Name)}, nil
}

func TestMain(m *testing.M) {
	var servers []server.Server

	if !testing.Short() {
		for i := 0; i < numServers; i++ {
			p := PORT + i
			s := server.NewServer(
				server.Name(fmt.Sprintf("test_gopher_%d", p)),
				server.Addr(fmt.Sprintf(":%d", p)),
				server.GRPCRegister(func(g *grpc.Server) {
					pb.RegisterGreeterServer(g, &greeter{})
				}))
			go func() {
				if err := s.Run(); err != nil {
					log.Fatal(err)
				}
			}()
			servers = append(servers, s)
		}
	}
	result := m.Run()
	if !testing.Short() {
		for _, s := range servers {
			s.Stop()
		}
	}
	os.Exit(result)
}

func InitConnectors() (cons []*connector) {
	for i := 0; i < numServers; i++ {
		p := PORT + i
		c := newConnector(
			Target(fmt.Sprintf("localhost:%d", p)),
			GRPCRegister(func(cc *grpc.ClientConn) interface{} {
				return pb.NewGreeterClient(cc)
			}))
		c.Init()
		cons = append(cons, c)
	}
	return cons
}

func InitBalancer(inCh <-chan Request) *Balancer {
	b := NewBalancer()
	go b.Balance(inCh)
	return b
}

func TestBalancer(t *testing.T) {

	requestsCh := make(chan Request)

	b := InitBalancer(requestsCh)
	conns := InitConnectors()
	// fill the pool with connectors
	for _, c := range conns {
		b.pool.Push(c)
	}
	t.Logf("Balancer %v", b)

	connsCh := make(chan *connector)
	// fake some requests
	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			r := Request{connsCh: connsCh}
			time.Sleep(time.Millisecond * 100)
			// send the call over the calls channel
			requestsCh <- r
		}()
	}

	// read connector from connsCh
	wg.Add(numRequests)
	for i := 0; i < numRequests; i++ {
		go func(i int) {
			defer wg.Done()
			conn := <-connsCh
			// Make a Call
			r, err := conn.Client().(pb.GreeterClient).SayHello(context.Background(), &pb.HelloRequest{Name: fmt.Sprintf("Gopher %d", i)})
			if err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, fmt.Sprintf("Hello Gopher %d", i), r.Message)
			// send conn over balancer connsCh
			b.connsCh <- conn
		}(i)
	}

	wg.Wait()
	// close connectors
	for _, c := range conns {
		c.Close()
	}
	log.Printf("TestBalancer END")
}

package server_test

import (
	"testing"
	"net/http"
	"github.com/paulormart/assert"
	"reflect"
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"fmt"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
	"bitbucket.org/aukbit/pluto/reply"
	pb "bitbucket.org/aukbit/pluto/server/proto"
	"io/ioutil"
	"encoding/json"
	"io"
)

func Home(w http.ResponseWriter, r *http.Request) {
  	reply.Json(w, r, http.StatusOK, "Hello World")
}

func Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reply.Json(w, r, http.StatusOK, fmt.Sprintf("Hello Room %s", ctx.Value("id").(string)))
}

type greeter struct{
	cfg 			*server.Config
}

// SayHello implements helloworld.GreeterServer
func (s *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: fmt.Sprintf("%v: Hello " + in.Name, s.cfg.Name)}, nil
}

func TestServer(t *testing.T){

	// HTTP server
	//1. create new server
	s := server.NewServer(
		server.Name("gopher"),
		server.Description("gopher super server"),
		server.Addr(":8080"),
	)
	assert.Equal(t, reflect.TypeOf(server.DefaultServer), reflect.TypeOf(s))

	cfg := s.Config()
	assert.Equal(t, true, len(cfg.Id) > 0)
	assert.Equal(t, "server_gopher", cfg.Name)
	assert.Equal(t, "gopher super server", cfg.Description)
	assert.Equal(t, ":8080", cfg.Addr)

	//2. register handlers
	mux := router.NewRouter()
	mux.GET("/home", Home)
	mux.GET("/home/:id", Detail)

	//3. assign last configs to the server before start, in this case setup a router
	s.Init(server.Mux(mux))

	//4. Run server
	go func(){
		if err := s.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	defer s.Stop()

	// GRPC server
	//1. create new server
	g := server.NewGRPCServer(
		server.Name("gopher"),
		server.Description("gopher super server"),
		server.Addr(":65057"),
	)

	cfg2 := g.Config()
	assert.Equal(t, true, len(cfg2.Id) > 0)
	assert.Equal(t, "server_gopher", cfg2.Name)
	assert.Equal(t, "gopher super server", cfg2.Description)

	// Register RegisterServerFunc
	g.Init(server.RegisterServerFunc(func(srv *grpc.Server) {
			pb.RegisterGreeterServer(srv, &greeter{cfg: cfg2})
		}),)

	// 2. Add some context
	go func() {
		//2. Run server
		if err := g.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	defer g.Stop()

	// Test
	const URL = "http://localhost:8080"
	var tests = []struct {
		Path         string
		Body         io.Reader
		BodyContains string
		Status       int
	}{
		{
		Path:         "/home",
		BodyContains: `Hello World`,
		Status:       http.StatusOK,
	},
		{
		Path:         "/home/123",
		BodyContains: `Hello Room 123`,
		Status:       http.StatusOK,
	},
	}
	for _, test := range tests {

		r, err := http.Get(URL + test.Path)
		if err != nil {
			log.Fatal(err)
		}
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer r.Body.Close()


		var message string
		if err := json.Unmarshal(b, &message); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, test.Status, r.StatusCode)
		assert.Equal(t, test.BodyContains, message)

	}




}


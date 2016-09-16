package server_test

import (
	"testing"
	"net/http"
	"pluto/server"
	"pluto/reply"
	"github.com/paulormart/assert"
	"reflect"
	"log"
	"pluto/server/router"
	"golang.org/x/net/context"
	pb "pluto/server/proto"
	"google.golang.org/grpc"
	"time"
	"fmt"
)

func Home(w http.ResponseWriter, r *http.Request) {
  	reply.Json(w, r, http.StatusOK, "Hello World")
}

func Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := map[string]string{"message": "Hello Room", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
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
	assert.Equal(t, "gopher.server.http", cfg.Name)
	assert.Equal(t, "gopher super server", cfg.Description)
	assert.Equal(t, ":8080", cfg.Addr)

	//2. register handlers
	mux := router.NewRouter()
	mux.GET("/home", Home)
	mux.GET("/home/:id", Detail)

	//3. assign last configs to the server before start, in this case setup a router
	s.Init(server.Router(mux))

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
	assert.Equal(t, "gopher.server.grpc", cfg2.Name)
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

	time.Sleep(time.Second * 600)

}


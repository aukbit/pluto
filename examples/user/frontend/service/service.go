package frontend

import (
	"flag"

	"github.com/aukbit/pluto/v6"
	"github.com/aukbit/pluto/v6/client"
	"github.com/aukbit/pluto/v6/examples/user/frontend/views"
	pb "github.com/aukbit/pluto/v6/examples/user/proto"
	"github.com/aukbit/pluto/v6/server"
	"github.com/aukbit/pluto/v6/server/router"
	"google.golang.org/grpc"
)

var target, httpPort string

func init() {
	flag.StringVar(&target, "target", "127.0.0.1:65087", "backend address")
	flag.StringVar(&httpPort, "http_port", ":8087", "frontend http port")
	flag.Parse()
}

func Run() error {
	// Define handlers
	mux := router.New()
	mux.Handle("GET", "/user", router.WrapErr(frontend.GetHandler))
	mux.Handle("GET", "/stream", router.WrapErr(frontend.GetStreamHandler))
	mux.Handle("POST", "/user", router.WrapErr(frontend.PostHandler))
	mux.Handle("GET", "/user/:id", router.WrapErr(frontend.GetHandlerDetail))
	mux.Handle("PUT", "/user/:id", router.WrapErr(frontend.PutHandler))
	mux.Handle("DELETE", "/user/:id", router.WrapErr(frontend.DeleteHandler))

	// define http server
	srv := server.New(
		server.Name("api"),
		server.Addr(httpPort),
		server.Mux(mux),
	)
	// Define grpc Client
	clt := client.New(
		client.Name("user"),
		client.GRPCRegister(func(cc *grpc.ClientConn) interface{} {
			return pb.NewUserServiceClient(cc)
		}),
		client.Target(target),
	)

	// Define Pluto service
	s := pluto.New(
		pluto.Name("frontend"),
		pluto.Description("Frontend service is responsible to parse all json data to regarding users to internal services"),
		pluto.Servers(srv),
		pluto.Clients(clt),
		pluto.HealthAddr(":9097"),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}

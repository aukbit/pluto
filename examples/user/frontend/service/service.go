package frontend

import (
	"flag"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/examples/user/frontend/views"
	pb "bitbucket.org/aukbit/pluto/examples/user/proto"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
	"google.golang.org/grpc"
)

var target = flag.String("target", "127.0.0.1:65060", "backend address")
var http_port = flag.String("http_port", ":8080", "frontend http port")

func Run() error {
	flag.Parse()

	// Define handlers
	mux := router.NewRouter()
	mux.GET("/user", frontend.GetHandler)
	mux.POST("/user", frontend.PostHandler)
	mux.GET("/user/:id", frontend.GetHandlerDetail)
	mux.PUT("/user/:id", frontend.PutHandler)
	mux.DELETE("/user/:id", frontend.DeleteHandler)

	// define http server
	httpSrv := server.NewServer(
		server.Name("api"),
		server.Addr(*http_port),
		server.Mux(mux))

	// Define grpc Client
	grpcClient := client.NewClient(
		client.Name("user"),
		client.RegisterClientFunc(func(cc *grpc.ClientConn) interface{} {
			return pb.NewUserServiceClient(cc)
		}),
		client.Target(*target),
	)
	// Define Pluto service
	s := pluto.NewService(
		pluto.Name("frontend"),
		pluto.Description("Frontend service is responsible to parse all json data to regarding users to internal services"),
	)
	// 5. Init service
	s.Init(pluto.Servers(httpSrv), pluto.Clients(grpcClient))

	// 6. Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}

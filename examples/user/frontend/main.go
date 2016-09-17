package main

import (
	"pluto"
	"pluto/server/router"
	"pluto/server"
	"pluto/client"
	"pluto/examples/user/frontend/views"
	pb "pluto/examples/user/proto"
	"google.golang.org/grpc"
	"context"
	"net/http"
	"log"
)

type user struct{}


func main(){

	// 1. Config service
	s := pluto.NewService(
		pluto.Name("frontend"),
		pluto.Description("user-frontend is responsible to parse all json data to regarding users to internal services"),
	)

	// 2. Set server handlers
	mux := router.NewRouter()

	//mux.GET("/user", WrapService(s, frontend.GetHandler))
	mux.GET("/user", frontend.GetHandler)
	mux.POST("/user", frontend.PostHandler)
	mux.GET("/user/:id", frontend.GetHandlerDetail)
	mux.PUT("/user/:id", frontend.PutHandler)
	mux.DELETE("/user/:id", frontend.DeleteHandler)

	// 3. Create new http server
	httpSrv := server.NewServer(server.Name("api"), server.Mux(mux))

	// 4. Define grpc Client
	grpcClient := client.NewClient(
		client.Name("user"),
		client.RegisterClientFunc(func(cc *grpc.ClientConn) interface{} {
			return pb.NewUserServiceClient(cc)
		}),
		client.Target("127.0.0.1:65059"),
	)
	// 5. Init service
	s.Init(pluto.Servers(httpSrv), pluto.Clients(grpcClient))

	// 6. Run service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

}

// TODO WrapService should live inside a service
func WrapService(s pluto.Service, next router.Handler) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "service", s)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
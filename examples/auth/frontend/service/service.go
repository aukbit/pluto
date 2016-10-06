package frontend

import (
	"flag"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/auth"
	"bitbucket.org/aukbit/pluto/examples/auth/frontend/views"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
)

var target = flag.String("target", "127.0.0.1:65081", "auth backend address")
var httpPort = flag.String("http_port", ":8081", "auth frontend http port")

// Run runs auth frontend service
func Run() error {
	flag.Parse()

	// Define handlers
	mux := router.NewMux()
	mux.POST("/authenticate", frontend.PostHandler)

	// define http server
	srv := server.NewServer(
		server.Name("api"),
		server.Addr(*httpPort),
		server.Mux(mux))

	// Define grpc Client
	clt := auth.NewClientAuth(*target)

	// Define Pluto service
	s := pluto.NewService(
		pluto.Name("auth_frontend"),
		pluto.Description("Authentication service is responsible to parse all json data to internal services"),
		pluto.Servers(srv),
		pluto.Clients(clt))

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}

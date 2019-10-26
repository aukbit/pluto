package frontend

import (
	"flag"

	"github.com/aukbit/pluto/v6"
	"github.com/aukbit/pluto/v6/auth"
	"github.com/aukbit/pluto/v6/examples/auth/frontend/views"
	"github.com/aukbit/pluto/v6/server"
	"github.com/aukbit/pluto/v6/server/router"
)

var (
	target   string
	httpPort string
)

func init() {
	flag.StringVar(&target, "target", "127.0.0.1:65081", "auth backend address")
	flag.StringVar(&httpPort, "http_port", ":8089", "auth frontend http port")
	flag.Parse()
}

// Run runs auth frontend service
func Run() error {

	// Define handlers
	mux := router.New()
	mux.POST("/authenticate", frontend.PostHandler)

	// define http server
	srv := server.New(
		server.Name("api"),
		server.Addr(httpPort),
		server.Mux(mux),
	)

	// Define grpc Client
	clt := auth.NewClientAuth(target)

	// Define Pluto service
	s := pluto.New(
		pluto.Name("auth_frontend"),
		pluto.Description("Authentication service is responsible to parse all json data to internal services"),
		pluto.Servers(srv),
		pluto.Clients(clt),
		pluto.HealthAddr(":9093"),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}

package web

import (
	"flag"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
)

var https_port = flag.String("https_port", ":8443", "https port")

func Run() error {

	// Set server handlers
	mux := router.NewMux()
	mux.GET("/", GetHandler)

	// Create new http server
	srv := server.NewServer(server.Name("api"),
		server.TLSConfig("server.crt", "private.key"),
		server.Addr(*https_port),
		server.Mux(mux))

	// Init service
	s := pluto.NewService(
		pluto.Name("web"),
		pluto.Description("web server serving handlers with https/tls"),
		pluto.Servers(srv),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}

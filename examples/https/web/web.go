package web

import (
	"flag"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
)

var https_port = flag.String("https_port", ":8443", "https port")

func Run() error {
	// 1. Config service
	s := pluto.NewService(
		pluto.Name("web"),
		pluto.Description("web server serving handlers with https/tls"),
	)
	// 2. Set server handlers
	mux := router.NewMux()
	mux.GET("/", GetHandler)
	// 3. Create new http server
	httpSrv := server.NewServer(server.Name("api"),
		server.TLSConfig("server.crt", "private.key"),
		server.Addr(*https_port),
		server.Mux(mux))
	// 4. Init service
	s.Init(pluto.Servers(httpSrv))

	// 5. Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}

package web

import (
	"log"
	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
	"flag"
)

func main() {
	log.Printf("teste")

}

var http_port = flag.String("http_port", ":443", "frontend http port")

func Run() error {
	log.SetFlags(log.Lshortfile)
	// 1. Config service
	s := pluto.NewService(
		pluto.Name("web"),
		pluto.Description("web server serving handlers with https/tls"),
	)
	// 2. Set server handlers
	mux := router.NewRouter()
	mux.GET("/", GetHandler)
	// 3. Create new http server
	httpSrv := server.NewServer(server.Name("api"),
		//server.TLSConfig("server.crt", "server.key"),
		server.Addr(*http_port),
		server.Mux(mux))
	// 4. Init service
	s.Init(pluto.Servers(httpSrv))

	// 5. Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}
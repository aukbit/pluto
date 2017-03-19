package main

import (
	"flag"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
)

var httpsPort = flag.String("https_port", ":8443", "https port")

func main() {
	// run frontend service
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Set server handlers
	mux := router.New()
	mux.GET("/", GetHandler)

	// Create new http server
	srv := server.New(server.Name("api"),
		server.TLSConfig("server.crt", "private.key"),
		server.Addr(*httpsPort),
		server.Mux(mux),
	)

	// Logger
	logger, _ := zap.NewDevelopment()
	// Init service
	s := pluto.New(
		pluto.Name("web"),
		pluto.Description("web server serving handlers with https/tls"),
		pluto.Servers(srv),
		pluto.Logger(logger),
		pluto.HealthAddr(":9098"),
	)

	// Run service
	if err := s.Run(); err != nil {
		return err
	}
	return nil
}

type Message struct {
	Message string `json:"message"`
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	m := &Message{"Hello Gopher"}
	reply.Json(w, r, http.StatusOK, m)
}

package main

import (
	"log"
	"net/http"

	"github.com/aukbit/pluto/v6"
	"github.com/aukbit/pluto/v6/reply"
	"github.com/aukbit/pluto/v6/server"
	"github.com/aukbit/pluto/v6/server/router"
)

func main() {
	// Define router
	mux := router.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		reply.Json(w, r, http.StatusOK, "Hello World")
	})

	// Define http server
	srv := server.New(
		server.Mux(mux),
	)

	// Define Pluto service
	s := pluto.New(
		pluto.Servers(srv),
	)

	// Run Pluto service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

package pluto

import (
	"log"
	"net/http"

	"bitbucket.org/aukbit/pluto/reply"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h := ctx.Value("pluto").(Service).Health()
	log.Printf("healthHandler %v", h)
	reply.Json(w, r, http.StatusOK, "Hello World")
}

func newHealthServer() server.Server {
	// Define Router
	mux := router.NewMux()
	mux.GET("/_health", healthHandler)
	// Define server
	return server.NewServer(server.Name("health"), server.Addr(":9090"), server.Mux(mux))
}

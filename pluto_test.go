package pluto_test

import (
	"testing"
	"net/http"
	"github.com/paulormart/assert"
	"pluto/reply"
	"pluto"
	"pluto/server/router"
	"log"
	"pluto/server"
)

func Index(w http.ResponseWriter, r *http.Request) {
  	reply.Json(w, r, http.StatusOK, "Hello World")
}

func TestServer(t *testing.T){

	// 1. Config service
	s := pluto.NewService(
		pluto.Name("gopher"),
		pluto.Description("gopher super service"),
	)
	//assert.Equal(t, reflect.TypeOf(service.DefaultServer), reflect.TypeOf(s))
	cfg := s.Config()
	assert.Equal(t, true, len(cfg.Id) > 0)
	assert.Equal(t, "gopher.pluto", cfg.Name)
	assert.Equal(t, "gopher super service", cfg.Description)
	assert.Equal(t, "gopher.pluto.server", cfg.Server.Config().Name)

	// 2. Set server handlers
	mux := router.NewRouter()
	mux.GET("/", Index)
	// 3. Define server Router
	s.Server().Init(server.Router(mux))

	// 4. Run service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

}
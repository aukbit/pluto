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

func TestService(t *testing.T){

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

	// 2. Set http server handlers
	mux := router.NewRouter()
	mux.GET("/", Index)
	// 3. Define server Router
	httpSrv := server.NewServer(server.Mux(mux))

	// 4. Init service
	s.Init(pluto.Servers(httpSrv))

	// 5. Run service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

}
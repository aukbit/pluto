package server_test

import (
	"testing"
	"net/http"
	"pluto/server"
	"pluto/reply"
	"github.com/paulormart/assert"
	"reflect"
	"log"
	"pluto/server/router"
)

func Home(w http.ResponseWriter, r *http.Request) {
  	reply.Json(w, r, http.StatusOK, "Hello World")
}

func Detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := map[string]string{"message": "Hello Room", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}

func TestServer(t *testing.T){

	//1. create new server
	s := server.NewServer(
		server.Name("gopher"),
		server.Description("gopher super server"),
		server.Addr(":8080"),
	)
	assert.Equal(t, reflect.TypeOf(server.DefaultServer), reflect.TypeOf(s))

	cfg := s.Config()
	assert.Equal(t, true, len(cfg.Id) > 0)
	assert.Equal(t, "gopher.server", cfg.Name)
	assert.Equal(t, "gopher super server", cfg.Description)
	assert.Equal(t, ":8080", cfg.Addr)

	//2. register handlers
	mux := router.NewRouter()
	mux.GET("/home", Home)
	mux.GET("/home/:id", Detail)

	//3. Define Router
	s.Router(mux)

	//4. Run server
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
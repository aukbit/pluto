package pluto_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"bitbucket.org/aukbit/pluto"
	"bitbucket.org/aukbit/pluto/reply"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
	"github.com/paulormart/assert"
)

func Index(w http.ResponseWriter, r *http.Request) {
	reply.Json(w, r, http.StatusOK, "Hello World")
}

func TestService(t *testing.T) {

	// Define Router
	mux := router.NewMux()
	mux.GET("/", Index)
	// Define server
	httpSrv := server.NewServer(server.Name("gopher"), server.Addr(":8083"), server.Mux(mux))

	// Define Service
	s := pluto.NewService(
		pluto.Name("gopher"),
		pluto.Description("gopher super service"),
		pluto.Servers(httpSrv),
	)

	// 5. Run service
	go func() {
		if err := s.Run(); err != nil {
			t.Fatal(err)
		}
	}()
	defer s.Stop()
	//
	time.Sleep(time.Second)

	// Assert Config
	cfg := s.Config()
	assert.Equal(t, true, len(cfg.ID) > 0)
	assert.Equal(t, "pluto_gopher", cfg.Name)
	assert.Equal(t, "gopher super service", cfg.Description)

	// Test
	const URL = "http://localhost:8083"
	r, err := http.Get(URL)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Body.Close()

	var message string
	if err := json.Unmarshal(b, &message); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, r.StatusCode)
	assert.Equal(t, "Hello World", message)

}

package router_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server/router"
	"github.com/paulormart/assert"
)

func Index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	reply.Json(w, r, http.StatusOK, "Hello World")
}
func GetHandler(w http.ResponseWriter, r *http.Request) {
	reply.Json(w, r, http.StatusOK, "Hello World")
}
func PostHandler(w http.ResponseWriter, r *http.Request) {
	var v interface{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer r.Body.Close()
	reply.Json(w, r, http.StatusCreated, v)
}
func GetDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := map[string]string{"message": "Hello World", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}
func PutDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := map[string]string{"message": "Hello World", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}
func DeleteDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := map[string]string{"message": "deleted", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}

func GetRoomHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := map[string]string{"message": "Hello World", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}

func GetCategoryDetailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data := map[string]string{"message": "Hello World",
		"id":       ctx.Value("id").(string),
		"category": ctx.Value("category").(string),
	}
	reply.Json(w, r, http.StatusOK, data)
}

func TestRouter(t *testing.T) {
	router := router.NewRouter()
	router.HandleFunc("GET", "/", IndexHandler)
	router.HandleFunc("GET", "/home", GetHandler)
	router.HandleFunc("GET", "/home/home", GetHandler)
	router.HandleFunc("GET", "/home/home/home", GetHandler)

	var tests = []struct {
		Method       string
		Path         string
		Body         io.Reader
		BodyContains interface{}
		Status       int
	}{
		{
			Method:       "GET",
			Path:         "/",
			BodyContains: "Hello World",
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home",
			BodyContains: "Hello World",
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/home",
			BodyContains: "Hello World",
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/home/home",
			BodyContains: "Hello World",
			Status:       http.StatusOK,
		},
	}
	server := httptest.NewServer(router)
	defer server.Close()
	for _, test := range tests {
		r, err := http.NewRequest(test.Method, server.URL+test.Path, test.Body)
		if err != nil {
			t.Fatal(err)
		}
		// call handler
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var v interface{}
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.Equal(t, test.Status, resp.StatusCode)
		assert.Equal(t, test.BodyContains, v)
	}
}

func TestDynamicRouter(t *testing.T) {
	router := router.NewRouter()
	router.HandleFunc("GET", "/", IndexHandler)
	router.HandleFunc("POST", "/home", PostHandler)
	router.HandleFunc("GET", "/home/:id", GetDetailHandler)
	router.HandleFunc("PUT", "/home/:id", PutDetailHandler)
	router.HandleFunc("DELETE", "/home/:id", DeleteDetailHandler)
	router.HandleFunc("GET", "/home/:id/room", GetRoomHandler)
	router.HandleFunc("GET", "/home/:id/room/:category", GetCategoryDetailHandler)

	var tests = []struct {
		Method       string
		Path         string
		Body         io.Reader
		BodyContains interface{}
		Status       int
	}{
		{
			Method:       "GET",
			Path:         "/",
			BodyContains: "Hello World",
			Status:       http.StatusOK,
		},
		{
			Method:       "POST",
			Path:         "/home",
			Body:         strings.NewReader(`{"name":"Gopher house"}`),
			BodyContains: map[string]string{"name": "Gopher house"},
			Status:       http.StatusCreated,
		},
		{
			Method:       "GET",
			Path:         "/home/123",
			BodyContains: map[string]string{"id": "123", "message": "Hello World"},
			Status:       http.StatusOK,
		},
		{
			Method: "PUT",
			Path:   "/home/123",

			Body:         strings.NewReader(`{"name":"Super Gopher house"}`),
			BodyContains: map[string]string{"id": "123", "message": "Hello World"},
			Status:       http.StatusOK,
		},
		{
			Method:       "DELETE",
			Path:         "/home/123",
			BodyContains: map[string]string{"id": "123", "message": "deleted"},
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/456/room",
			BodyContains: map[string]string{"id": "456", "message": "Hello World"},
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/456/room/999",
			BodyContains: map[string]string{"id": "456", "category": "999", "message": "Hello World"},
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/",
			BodyContains: "404 page not found",
			Status:       http.StatusNotFound,
		},
		{
			Method:       "GET",
			Path:         "/abc",
			BodyContains: "404 page not found",
			Status:       http.StatusNotFound,
		},
		{
			Method:       "GET",
			Path:         "/somethingelse/123/w444/f444",
			BodyContains: "404 page not found",
			Status:       http.StatusNotFound,
		},
	}
	server := httptest.NewServer(router)
	defer server.Close()
	for _, test := range tests {
		r, err := http.NewRequest(test.Method, server.URL+test.Path, test.Body)
		if err != nil {
			t.Fatal(err)
		}
		// call handler
		resp, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		var v interface{}
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.Equal(t, test.Status, resp.StatusCode)
		switch v.(type) {
		case string:
			assert.Equal(t, test.BodyContains, v)
		case map[string]interface{}:
			assert.Equal(t, test.BodyContains.(map[string]string)["name"], v.(map[string]interface{})["name"])
			assert.Equal(t, test.BodyContains.(map[string]string)["id"], v.(map[string]interface{})["id"])
			assert.Equal(t, test.BodyContains.(map[string]string)["category"], v.(map[string]interface{})["category"])
			assert.Equal(t, test.BodyContains.(map[string]string)["message"], v.(map[string]interface{})["message"])
		}
	}
}

// GOMAXPROCS=1 go test ./server/router -bench=BenchmarkRouter -benchmem
// BenchmarkRouter             1000           1945098 ns/op           21038 B/op        211 allocs/op
func BenchmarkRouter(b *testing.B) {
	router := router.NewRouter()
	router.HandleFunc("GET", "/", IndexHandler)
	router.HandleFunc("POST", "/home", PostHandler)
	router.HandleFunc("GET", "/home/:id", GetDetailHandler)
	router.HandleFunc("PUT", "/home/:id", PutDetailHandler)
	router.HandleFunc("DELETE", "/home/:id", DeleteDetailHandler)
	router.HandleFunc("GET", "/home/:id/room", GetRoomHandler)
	router.HandleFunc("GET", "/home/:id/room/:category", GetCategoryDetailHandler)

	server := httptest.NewServer(router)
	defer server.Close()
	// run the Put function b.N times
	for n := 0; n < b.N; n++ {
		// create new request
		r, err := http.NewRequest("GET", server.URL+"/home/456/room/999", strings.NewReader(`{}`))
		if err != nil {
			b.Fatal(err)
		}
		// call handler
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			b.Fatal(err)
		}
		response.Body.Close()
	}
}

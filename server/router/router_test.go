package router_test

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/paulormart/assert"

	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server/router"
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
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := make(map[string]string)
	if err := json.Unmarshal(body, &data); err != nil {
		log.Println(err.Error())
	}

	reply.Json(w, r, http.StatusCreated, data)
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

func TestServer(t *testing.T) {
	router := router.NewRouter()
	router.Handle("GET", "/", IndexHandler)
	router.Handle("GET", "/home", GetHandler)
	router.Handle("GET", "/home/home", GetHandler)
	router.Handle("GET", "/home/home/home", GetHandler)
	router.Handle("POST", "/home", PostHandler)
	router.Handle("GET", "/home/:id", GetDetailHandler)
	router.Handle("PUT", "/home/:id", PutDetailHandler)
	router.Handle("DELETE", "/home/:id", DeleteDetailHandler)
	router.Handle("GET", "/home/:id/room", GetRoomHandler)
	router.Handle("GET", "/home/:id/room/:category", GetCategoryDetailHandler)
	router.GET("/home", GetHandler)
	router.POST("/home", PostHandler)
	router.PUT("/home/:id", PutDetailHandler)
	router.DELETE("/home/:id", DeleteDetailHandler)

	var tests = []struct {
		Method       string
		Path         string
		Body         io.Reader
		BodyContains string
		Status       int
	}{
		{
			Method:       "GET",
			Path:         "/",
			BodyContains: `"Hello World"`,
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home",
			BodyContains: `"Hello World"`,
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/home",
			BodyContains: `"Hello World"`,
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/home/home",
			BodyContains: `"Hello World"`,
			Status:       http.StatusOK,
		},
		{
			Method:       "POST",
			Path:         "/home",
			Body:         strings.NewReader(`{"name":"Gopher house"}`),
			BodyContains: `{"name":"Gopher house"}`,
			Status:       http.StatusCreated,
		},
		{
			Method:       "GET",
			Path:         "/home/123",
			BodyContains: `{"id":"123","message":"Hello World"}`,
			Status:       http.StatusOK,
		},
		{
			Method:       "PUT",
			Path:         "/home/123",
			Body:         strings.NewReader(`{"name":"Super Gopher house"}`),
			BodyContains: `{"id":"123","message":"Hello World"}`,
			Status:       http.StatusOK,
		},
		{
			Method:       "DELETE",
			Path:         "/home/123",
			BodyContains: `{"id":"123","message":"deleted"}`,
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/456/room",
			BodyContains: `{"id":"456","message":"Hello World"}`,
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/456/room/999",
			BodyContains: `{"category":"999","id":"456","message":"Hello World"}`,
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/home/",
			BodyContains: `"Hello World"`,
			Status:       http.StatusOK,
		},
		{
			Method:       "GET",
			Path:         "/abc",
			BodyContains: `"404 page not found"`,
			Status:       http.StatusNotFound,
		},
		{
			Method:       "GET",
			Path:         "/somethingelse/123/w444/f444",
			BodyContains: `"404 page not found"`,
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
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		actualBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, response.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, test.Status, response.StatusCode)
		assert.Equal(t, test.BodyContains, string(actualBody))

	}

}

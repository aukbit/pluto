package router

import (
	"testing"
	"github.com/paulormart/assert"
	"pluto/reply"

	"io/ioutil"
	"net/http/httptest"
	"io"
	"net/http"
	"encoding/json"
	"log"
	"strings"
)

func TestPaths(t *testing.T){

	var path, key, value, prefix string
	var params []string
	path = "/"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/", key)
	assert.Equal(t, "/", value)
	assert.Equal(t, "", prefix)
	assert.Equal(t, []string{}, params)
	path = "/home"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home", key)
	assert.Equal(t, "/home", value)
	assert.Equal(t, "", prefix)
	assert.Equal(t, []string{}, params)
	path = "/home/:id"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home/:", key)
	assert.Equal(t, "/:", value)
	assert.Equal(t, "/home", prefix)
	assert.Equal(t, []string{"id"}, params)
	path = "/home/:id/room"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home/:/room", key)
	assert.Equal(t, "/room", value)
	assert.Equal(t, "/home/:", prefix)
	assert.Equal(t, []string{"id"}, params)

}

func Index (w http.ResponseWriter, r *http.Request){
	io.WriteString(w, "Hello World")
}

func TestHandlers(t *testing.T){

	router := NewRouter()
	router.Handle("GET", "/home", Index)
	assert.Equal(t, 1, router.trie.Size())
	assert.Equal(t, true, router.trie.Contains("/home"))

	data := router.trie.Get("/home")
	assert.Equal(t, "/home", data.value)
	assert.Equal(t, []string{}, data.vars)
	assert.Equal(t, 1, len(data.methods))
	assert.Equal(t, true, data.methods["GET"] != nil)

	router.Handle("POST", "/home", Index)
	assert.Equal(t, 1, router.trie.Size())
	assert.Equal(t, 2, len(data.methods))
	assert.Equal(t, true, data.methods["GET"] != nil)
	assert.Equal(t, true, data.methods["POST"] != nil)

	router.Handle("GET", "/home/:id", Index)
	router.Handle("PUT", "/home/:id", Index)
	router.Handle("DELETE", "/home/:id", Index)
	data1 := router.trie.Get("/home/:")
	assert.Equal(t, 2, router.trie.Size())
	assert.Equal(t, 3, len(data1.methods))
	assert.Equal(t, true, data1.methods["GET"] != nil)
	assert.Equal(t, true, data1.methods["PUT"] != nil)
	assert.Equal(t, true, data1.methods["DELETE"] != nil)

}

func IndexHandler (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}
func GetHandler (w http.ResponseWriter, r *http.Request){
	reply.Json(w, r, http.StatusOK, "Hello World")
}
func PostHandler (w http.ResponseWriter, r *http.Request){
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
func GetDetailHandler (w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	data := map[string]string{"message": "Hello World", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}
func PutDetailHandler (w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	data := map[string]string{"message": "Hello World", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}
func DeleteDetailHandler (w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	data := map[string]string{"message": "deleted", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}

func GetRoomHandler (w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	data := map[string]string{"message": "Hello World", "id": ctx.Value("id").(string)}
	reply.Json(w, r, http.StatusOK, data)
}

func GetCategoryDetailHandler (w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	data := map[string]string{"message": "Hello World",
		"id": ctx.Value("id").(string),
		"category": ctx.Value("category").(string),
	}
	reply.Json(w, r, http.StatusOK, data)
}

func TestServer(t *testing.T){
	router := NewRouter()
	router.Handle("GET", "/", IndexHandler)
	router.Handle("GET", "/home", GetHandler)
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
		r, err := http.NewRequest(test.Method, server.URL + test.Path, test.Body)
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
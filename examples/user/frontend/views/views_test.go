package frontend_test

import (
	"testing"
	"io"
	"strings"
	"net/http"
	"pluto/server/router"
	"net/http/httptest"
	"io/ioutil"
	"github.com/paulormart/assert"
	"pluto/examples/user/frontend/views"
)

func TestServer(t *testing.T){
	router := router.NewRouter()
	//router.GET("/user", GetHandler)
	router.POST("/user", frontend.PostHandler)
	//router.PUT("/home/:id", PutDetailHandler)
	//router.DELETE("/home/:id", DeleteDetailHandler)

	var tests = []struct {
		Method       string
		Path         string
		Body         io.Reader
		BodyContains string
		Status       int
	}{
	//	{
	//	Method:       "GET",
	//	Path:         "/user",
	//	BodyContains: `"Hello World"`,
	//	Status:       http.StatusOK,
	//},
		{
		Method:       "POST",
		Path:         "/user",
		Body:         strings.NewReader(`{"name":"Gopher", "email": "gopher@email.com", "password":"123456"}`),
		BodyContains: `{"name":"Gopher", "email": "gopher@email.com"}`,
		Status:       http.StatusCreated,
	},
	//	{
	//	Method:       "GET",
	//	Path:         "/user/123",
	//	BodyContains: `{"id":"123","message":"Hello World"}`,
	//	Status:       http.StatusOK,
	//},
	//	{
	//	Method:       "PUT",
	//	Path:         "/user/123",
	//	Body:         strings.NewReader(`{"name":"Super Gopher house"}`),
	//	BodyContains: `{"id":"123","message":"Hello World"}`,
	//	Status:       http.StatusOK,
	//},
	//	{
	//	Method:       "DELETE",
	//	Path:         "/home/user",
	//	BodyContains: `{"id":"123","message":"deleted"}`,
	//	Status:       http.StatusOK,
	//},
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
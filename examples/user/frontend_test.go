package frontend_test

import (
	"github.com/paulormart/assert"
	"pluto/examples/user/frontend/service"
	"testing"
	"log"
	"io/ioutil"
	"strings"
	"io"
	"net/http"
)

func TestAll(t *testing.T){

	// Note: Run the backend service in a terminal window

	// launch frontend service running on
	// default http://localhost:8080
	const URL = "http://localhost:8080"

	go func(){
		if err := frontend.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	//
	var tests = []struct {
		Method       string
		Path         string
		Body         io.Reader
		BodyContains string
		Status       int
	}{
		{
		Method:       "GET",
		Path:         "/user",
		BodyContains: `"Hello World"`,
		Status:       http.StatusOK,
	},
		{
		Method:       "POST",
		Path:         "/user",
		Body:         strings.NewReader(`{"name":"Gopher", "email": "gopher@email.com", "password":"123456"}`),
		BodyContains: `{"id":"123","name":"Gopher","email":"gopher@email.com"}`,
		Status:       http.StatusCreated,
	},
		{
		Method:       "GET",
		Path:         "/user/123",
		BodyContains: `{"id":"123","message":"Hello World"}`,
		Status:       http.StatusOK,
	},
		{
		Method:       "PUT",
		Path:         "/user/123",
		Body:         strings.NewReader(`{"name":"Super Gopher house"}`),
		BodyContains: `{"id":"123","message":"Hello World"}`,
		Status:       http.StatusOK,
	},
		{
		Method:       "DELETE",
		Path:         "/home/user",
		BodyContains: `{"id":"123","message":"deleted"}`,
		Status:       http.StatusOK,
	},
	}

	for _, test := range tests {
		r, err := http.NewRequest(test.Method, URL + test.Path, test.Body)
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

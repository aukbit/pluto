package frontend_test

import (
	"github.com/paulormart/assert"
	"bitbucket.org/aukbit/pluto/examples/user/frontend/service"
	"testing"
	"log"
	"io"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
)

type User struct {
	Id  	string           `json:"id"`
	Name  	string           `json:"name"`
	Email  	string           `json:"email"`
}

const URL = "http://localhost:8080"

func TestAll(t *testing.T){

	// Note: Run the backend service in a different terminal window
	// $  go run ./backend/main.go -db_addr=DB_ADDR_ENV

	// launch frontend service running on
	// default http://localhost:8080

	go func(){
		if err := frontend.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	//
	user := &User{}
	//
	var tests = []struct {
		Method       string
		Path         func(string)string
		Body         io.Reader
		BodyContains func(string)string
		Status       int
	}{
		{
		Method:       "POST",
		Path:         func(id string) string { return URL + "/user" },
		Body:         strings.NewReader(`{"name":"Gopher", "email": "gopher@email.com", "password":"123456"}`),
		BodyContains: func(id string) string { return `{"id":"`+id+`","name":"Gopher","email":"gopher@email.com"}` },
		Status:       http.StatusCreated,
	},
		{
		Method:       "GET",
		Path:         func(id string) string { return URL + "/user/" + id },
		BodyContains: func(id string) string { return `{"id":"`+id+`","name":"Gopher","email":"gopher@email.com"}` },
		Status:       http.StatusOK,
	},
		{
		Method:       "PUT",
		Path:         func(id string) string { return URL + "/user/" + id },
		Body:         strings.NewReader(`{"name":"Super Gopher house"}`),
		BodyContains: func(id string) string { return `{"id":"`+id+`","name":"Super Gopher house"}` },
		Status:       http.StatusOK,
	},
		{
		Method:       "DELETE",
		Path:         func(id string) string { return URL + "/user/" + id },
		BodyContains: func(id string) string { return `{}` },
		Status:       http.StatusOK,
	},
	}

	for _, test := range tests {

		r, err := http.NewRequest(test.Method, test.Path(user.Id), test.Body)
		if err != nil {
			t.Fatal(err)
		}
		// call handler
		response, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		actualBody, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		err = json.Unmarshal(actualBody, user)
		if err != nil {
			t.Fatalf("Unmarshal %v %v", string(actualBody), err)
		}
		assert.Equal(t, response.Header.Get("Content-Type"), "application/json")
		assert.Equal(t, test.Status, response.StatusCode)
		assert.Equal(t, test.BodyContains(user.Id), string(actualBody))
	}

	// FilterUsers
	r, err := http.NewRequest("GET", URL + "/user?name=Gopher", nil)
	if err != nil {
		t.Fatal(err)
	}
	// call handler
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	actualBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, response.Header.Get("Content-Type"), "application/json")
	assert.Equal(t, response.StatusCode, response.StatusCode)
	assert.Equal(t, true, len(actualBody) > 0)

}

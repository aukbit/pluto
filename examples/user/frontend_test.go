package frontend_test

import (
	"github.com/paulormart/assert"
	"pluto/examples/user/frontend/service"
	"testing"
	"log"
	"strings"
	"io"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type User struct {
	Id  	string           `json:"id"`
	Name  	string           `json:"name"`
	Email  	string           `json:"email"`
}

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
	user := &User{}
	//
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
	//	Status:       http.StatusOK,
	//},
		{
		Method:       "POST",
		Path:         "/user",
		Body:         strings.NewReader(`{"name":"Gopher", "email": "gopher@email.com", "password":"123456"}`),
		Status:       http.StatusCreated,
	},
		{
		Method:       "GET",
		Path:         "/user/",
		Status:       http.StatusOK,
	},
	//	{
	//	Method:       "PUT",
	//	Path:         "/user/" + user.Id,
	//	Body:         strings.NewReader(`{"name":"Super Gopher house"}`),
	//	Status:       http.StatusOK,
	//},
	//	{
	//	Method:       "DELETE",
	//	Path:         "/user/" + user.Id,
	//	Status:       http.StatusOK,
	//},
	}


	for _, test := range tests {

		var url string
		if user == nil {
			url = URL + test.Path
		} else {
			url = URL + test.Path + user.Id
		}

		r, err := http.NewRequest(test.Method, url, test.Body)
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
		assert.Equal(t, true, len(user.Id) > 0)
		assert.Equal(t, true, len(user.Email) > 0)
		assert.Equal(t, "Gopher", user.Name)
	}
}

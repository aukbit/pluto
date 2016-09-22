package web_test

import (
	"github.com/paulormart/assert"
	"bitbucket.org/aukbit/pluto/examples/https/web"
	"testing"
	"log"
	"io"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Message struct {
	message  string           `json:"message"`
}

const URL = "https://localhost:8443"

func TestAll(t *testing.T){

	// launch frontend service running on
	// default http://localhost:8080

	go func(){
		if err := web.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	var tests = []struct {
		Method       string
		Path         func()string
		Body         io.Reader
		BodyContains func()string
		Status       int
	}{
		{
		Method:       "GET",
		Path:         func() string { return URL + "/" },
		BodyContains: func() string { return `{"message":"Hello Gopher"}` },
		Status:       http.StatusCreated,
	},

	}

	message := &Message{}

	for _, test := range tests {

		r, err := http.NewRequest(test.Method, test.Path(), test.Body)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, 1, 1)
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
		err = json.Unmarshal(actualBody, message)
		if err != nil {
			log.Fatalf("Unmarshal %v", err)
		} else {
			assert.Equal(t, response.Header.Get("Content-Type"), "application/json")
			assert.Equal(t, test.Status, response.StatusCode)
			assert.Equal(t, test.BodyContains(), string(actualBody))
		}

	}

}

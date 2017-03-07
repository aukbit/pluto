package web_test

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/aukbit/pluto/examples/https/web"
	"github.com/paulormart/assert"
)

const URL = "https://localhost:8443"

func TestMain(m *testing.M) {

	if !testing.Short() {
		// Run Server
		go func() {
			if err := web.Run(); err != nil {
				log.Fatal(err)
			}
		}()
		time.Sleep(time.Second)
	}
	result := m.Run()
	if !testing.Short() {
		// Stop Server
	}
	os.Exit(result)
}

func TestExampleHTTPS(t *testing.T) {

	// Note: launch frontend service in a terminal rather than lunching a go routine here
	// $ go run main.go
	// default http://localhost:8080

	//go func(){
	//	if err := web.Run(); err != nil {
	//		log.Fatal(err)
	//	}
	//}()

	var tests = []struct {
		Method       string
		Path         func() string
		Body         io.Reader
		BodyContains func() string
		Status       int
	}{
		{
			Method:       "GET",
			Path:         func() string { return URL + "/" },
			BodyContains: func() string { return `{"message":"Hello Gopher"}` },
			Status:       http.StatusOK,
		},
	}

	message := &web.Message{}

	for _, test := range tests {

		// skip certificate verification
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		response, err := client.Get(test.Path())
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
			assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
			assert.Equal(t, "max-age=63072000; includeSubDomains", response.Header.Get("Strict-Transport-Security"))
			assert.Equal(t, test.Status, response.StatusCode)
			assert.Equal(t, test.BodyContains(), string(actualBody))
		}

	}

}

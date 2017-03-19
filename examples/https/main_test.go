package main

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/paulormart/assert"
)

const URL = "https://localhost:8443"

func TestMain(m *testing.M) {

	if !testing.Short() {
		// Run Server
		go func() {
			if err := run(); err != nil {
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

	var tests = []struct {
		Method       string
		Path         func() string
		Body         io.Reader
		BodyContains func() *Message
		Status       int
	}{
		{
			Method:       "GET",
			Path:         func() string { return URL + "/" },
			BodyContains: func() *Message { return &Message{"Hello Gopher"} },
			Status:       http.StatusOK,
		},
	}

	message := &Message{}

	for _, test := range tests {

		// skip certificate verification
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(test.Path())
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if err = json.NewDecoder(resp.Body).Decode(&message); err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
		assert.Equal(t, "max-age=63072000; includeSubDomains", resp.Header.Get("Strict-Transport-Security"))
		assert.Equal(t, test.Status, resp.StatusCode)
		assert.Equal(t, test.BodyContains(), message)
	}
}

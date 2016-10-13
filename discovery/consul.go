package discovery

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var URL = "http://192.168.99.100:8500"

const (
	SELF = "/v1/agent/self" // Returns the local node configuration
)

func IsAvailable() (bool, error) {
	resp, err := http.Get(URL + SELF)
	if err != nil {
		return false, err
	}
	_, err = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Consul not available at this node")
	}
	return true, nil
}

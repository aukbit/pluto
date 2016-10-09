package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var URL = "http://192.168.99.100:8500"

const (
	register    = "/v1/catalog/register"
	deregister  = "/v1/catalog/deregister"
	datacenters = "/v1/catalog/datacenters"
	nodes       = "/v1/catalog/nodes"
	services    = "/v1/catalog/services"
)

type consul struct{}

func newT() {

}

type R struct {
	Nodes []Node
}

type Node struct {
	Node    string `json:"node"`
	Address string `json:"address"`
	// TaggedAddresses  string `json:"TaggedAddresses"`
}

func Nodes() {

	qs := "?near=_agent"
	resp, err := http.Get(URL + nodes + qs)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("Nodes %v", err)
		// return err
	}
	// nds := &R{}
	var data []json.RawMessage
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range data {
		node := &Node{}
		if err := json.Unmarshal(n, node); err != nil {
			fmt.Println(err)
		}
		log.Printf("Nodes %v", node)
	}

	// log.Printf("Nodes %v", data)

}

func Register() error {
	req, err := http.NewRequest("PUT", URL+register, strings.NewReader(`{"Service": {"ID":"test1", "Service": "redis"}}`))
	if err != nil {
		log.Printf("Register %v", err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Register %v", err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("Register %v", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		e := errors.New(fmt.Sprintf("Error %v", string(body)))
		return e
	}
	return nil
}

func Unregister() {

}

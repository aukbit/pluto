package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	catalogNodesPath = "/v1/catalog/nodes" //Lists nodes in a given DC
)

// Node single consul node
type Node struct {
	Node            string            `json:"Node"`
	Address         string            `json:"Address"`
	TaggedAddresses map[string]string `json:"TaggedAddresses,omitempty"`
}

// Nodes slice of nodes
type Nodes []Node

// Noder interface
type Noder interface {
	GetNodes(addr, path string) (Nodes, error)
}

// DefaultNoder struct to append GetNodes
type DefaultNoder struct{}

// GetNodes make GET request on consul api
func (dn *DefaultNoder) GetNodes(addr, path string) (Nodes, error) {
	qs := "?near=_agent"
	url := fmt.Sprintf("http://%s%s%s", addr, path, qs)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	nodes := Nodes{}
	if err := json.Unmarshal(body, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetNodes function to get slice of nodes
func GetNodes(n Noder, addr string) (Nodes, error) {
	nodes, err := n.GetNodes(addr, catalogNodesPath)
	if err != nil {
		return nil, fmt.Errorf("Error querying Consul API: %s", err)
	}
	return nodes, nil
}

package discovery

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	NODES = "/v1/catalog/nodes"
)

type Node struct {
	Node            string            `json:"Node"`
	Address         string            `json:"Address"`
	TaggedAddresses map[string]string `json:"TaggedAddresses"`
}

func Nodes() (nodes []*Node, err error) {

	qs := "?near=_agent"
	resp, err := http.Get(URL + NODES + qs)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var data []json.RawMessage
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	for _, node := range data {
		n := &Node{}
		if err := json.Unmarshal(node, n); err != nil {
			return nil, err
		}
		log.Printf("nodes %v", n)
		nodes = append(nodes, n)
	}
	return nodes, nil
}

package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	catalogServicePath = "/v1/catalog/service" // Lists the nodes in a given service
)

// NodeService single consul node-service
type NodeService struct {
	Node           string   `json:"Node"`
	Address        string   `json:"Address"`
	ServiceID      string   `json:"ServiceID"`
	ServiceName    string   `json:"ServiceName"`
	ServiceTags    []string `json:"ServiceTags,omitempty"`
	ServiceAddress string   `json:"ServiceAddress"`
	ServicePort    int      `json:"ServicePort"`
}

// NodeServices slice of node-services
type NodeServices []NodeService

// NodeServicer interface
type NodeServicer interface {
	GetNodeServices(addr, path string) (NodeServices, error)
}

// GetNodeServices function to get slice of node-services
func GetNodeServices(n NodeServicer, addr string) (NodeServices, error) {
	nodes, err := n.GetNodeServices(addr, catalogServicePath)
	if err != nil {
		return nil, fmt.Errorf("Error querying Consul API: %s", err)
	}
	return nodes, nil
}

// DefaultNodeServicer struct to append GetNodeServices
type DefaultNodeServicer struct{}

// GetNodeServices make GET request on consul api
func (dn *DefaultNodeServicer) GetNodeServices(addr, path string) (NodeServices, error) {
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
	nodes := NodeServices{}
	if err := json.Unmarshal(body, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

// func CatalogService(url, service string) (ns []*NodeService, err error) {
// 	if service == "" {
// 		return nil, fmt.Errorf("to search for a service in service discovery, a service name must be specified")
// 	}
// 	resp, err := http.Get(url + strings.Replace(catalogServicePath, "<service>", service, 1))
// 	if err != nil {
// 		return nil, err
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data []json.RawMessage
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, bytes := range data {
// 		n := &NodeService{}
// 		err = json.Unmarshal(bytes, n)
// 		if err != nil {
// 			return nil, err
// 		}
// 		ns = append(ns, n)
// 	}
// 	return ns, nil
// }

// Targets create a slice of addresses based on the services
// returned from CatalogService
// func Targets(url, service string) (targets []string, err error) {
// 	ns, err := CatalogService(url, service)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(ns) == 0 {
// 		return nil, fmt.Errorf("nodes not available with service: %s", service)
// 	}
// 	for _, n := range ns {
// 		t := fmt.Sprintf("%s:%d", n.Address, n.ServicePort)
// 		targets = append(targets, t)
// 	}
// 	return targets, nil
// }

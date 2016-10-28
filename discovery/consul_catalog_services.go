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

// ServiceNode single consul service-node
type ServiceNode struct {
	Node           string   `json:"Node"`
	Address        string   `json:"Address"`
	ServiceID      string   `json:"ServiceID"`
	ServiceName    string   `json:"ServiceName"`
	ServiceTags    []string `json:"ServiceTags,omitempty"`
	ServiceAddress string   `json:"ServiceAddress"`
	ServicePort    int      `json:"ServicePort"`
}

// ServiceNodes slice of service-nodes
type ServiceNodes []ServiceNode

// ServiceNoder interface
type ServiceNoder interface {
	GetServiceNodes(addr, path string, serviceID string) (ServiceNodes, error)
}

// GetServiceNodes function to get slice of service-nodes
func GetServiceNodes(s ServiceNoder, addr, serviceID string) (ServiceNodes, error) {
	nodes, err := s.GetServiceNodes(addr, catalogServicePath, serviceID)
	if err != nil {
		return nil, fmt.Errorf("Error querying Consul API: %s", err)
	}
	return nodes, nil
}

// DefaultServiceNoder struct to append GetNodeServices
type DefaultServiceNoder struct{}

// GetServiceNodes make GET request on consul api
func (dn *DefaultServiceNoder) GetServiceNodes(addr, path, serviceID string) (ServiceNodes, error) {
	qs := "?near=_agent"
	url := fmt.Sprintf("http://%s%s/%s%s", addr, path, serviceID, qs)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	nodes := ServiceNodes{}
	if err := json.Unmarshal(body, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

// GetServiceTargets create a slice of addresses based on the services
// returned from CatalogService
func GetServiceTargets(addr, serviceID string) (targets []string, err error) {
	ns, err := GetServiceNodes(&DefaultServiceNoder{}, addr, serviceID)
	if err != nil {
		return nil, err
	}
	if len(ns) == 0 {
		return nil, fmt.Errorf("Error service: %s is not available in any of the nodes", serviceID)
	}
	for _, n := range ns {
		t := fmt.Sprintf("%s:%d", n.Address, n.ServicePort)
		targets = append(targets, t)
	}
	return targets, nil
}

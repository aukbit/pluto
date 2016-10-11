package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	CATALOG_SERVICES = "/v1/catalog/services"          // Lists services in a given DC
	CATALOG_SERVICE  = "/v1/catalog/service/<service>" // Lists the nodes in a given service
)

type NodeService struct {
	Node           string   `json:"Node"`
	Address        string   `json:"Address"`
	ServiceID      string   `json:"ServiceID"`
	ServiceName    string   `json:"ServiceName"`
	ServiceTags    []string `json:"ServiceTags"`
	ServiceAddress string   `json:"ServiceAddress"`
	ServicePort    int      `json:"ServicePort"`
}

func CatalogServices() (map[string][]string, error) {

	resp, err := http.Get(URL + CATALOG_SERVICES)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var data map[string][]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func CatalogService(service string) (ns []*NodeService, err error) {

	resp, err := http.Get(URL + strings.Replace(CATALOG_SERVICE, "<service>", service, 1))
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
	for _, bytes := range data {
		n := &NodeService{}
		err = json.Unmarshal(bytes, n)
		if err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func Target(service string) (string, error) {
	ns, err := CatalogService(service)
	if err != nil {
		return "", err
	}
	if len(ns) > 0 {
		t := fmt.Sprintf("%s:%d", ns[0].Address, ns[0].ServicePort)
		return t, nil
	}
	return "", fmt.Errorf("nodes not available with service: %s", service)
}

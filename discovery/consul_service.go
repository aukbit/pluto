package discovery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	SERVICES           = "/v1/agent/services"                       // Returns the services the local agent is managing
	REGISTER_SERVICE   = "/v1/agent/service/register"               // Registers a new local service
	DEREGISTER_SERVICE = "/v1/agent/service/deregister/<serviceID>" // Deregisters a local service
)

type Service struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Tags    []string `json:"Tags,omitempty"`
	Address string   `json:"Address,omitempty"`
	Port    int      `json:"Port"`
}

func services(url string) (map[string]*Service, error) {

	resp, err := http.Get(url + SERVICES)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var data map[string]*Service
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func registerService(url string, s *Service) error {

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url+REGISTER_SERVICE, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error %v", string(body))
	}
	return nil
}

func deregisterService(url string, serviceID string) error {

	req, err := http.NewRequest("PUT", url+strings.Replace(DEREGISTER_SERVICE, "<serviceID>", serviceID, 1), bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error %v", string(body))
	}
	return nil
}

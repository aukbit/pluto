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
	REGISTER_CHECK   = "/v1/agent/check/register"             // Registers a new local check
	DEREGISTER_CHECK = "/v1/agent/check/deregister/<checkID>" // Deregisters a local check
)

type Check struct {
	ID                             string `json:"ID"`
	Name                           string `json:"Name"`
	Notes                          string `json:"Notes,omitempty"`
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter,omitempty"`
	Script                         string `json:"Script,omitempty"`
	DockerContainerID              string `json:"DockerContainerID,omitempty"`
	Shell                          string `json:"Shell,omitempty"`
	HTTP                           string `json:"HTTP,omitempty"`
	TCP                            string `json:"TCP,omitempty"`
	Interval                       string `json:"Interval,omitempty"`
	TTL                            string `json:"TTL,omitempty"`
	Timeout                        string `json:"Timeout,omitempty"`
	ServiceID                      string `json:"ServiceID,omitempty"`
}

func registerCheck(url string, c *Check) error {

	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", url+REGISTER_CHECK, bytes.NewBuffer(b))
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

func deregisterCheck(url string, checkID string) error {

	req, err := http.NewRequest("PUT", url+strings.Replace(DEREGISTER_CHECK, "<checkID>", checkID, 1), bytes.NewBuffer([]byte(`{}`)))
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

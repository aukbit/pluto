package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	checksPath          = "/v1/agent/checks"                     // Returns the checks the local agent is managing
	registerCheckPath   = "/v1/agent/check/register"             // Registers a new local check
	deregisterCheckPath = "/v1/agent/check/deregister/<checkID>" // Deregisters a local check
)

// Check struct
type Check struct {
	ID                             string `json:"ID,omitempty"`
	CheckID                        string `json:"CheckID,omitempty"`
	Name                           string `json:"Name"`
	Node                           string `json:"Node,omitempty"`
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
	Status                         string `json:"Status,omitempty"`
}

// Checks map of checks
type Checks map[string]Check

// Checker interface
type Checker interface {
	GetChecks(addr, path string) (Checks, error)
}

// DefaultServicer struct to implement Servicer default methods
type DefaultChecker struct{}

// GetChecks make GET request on consul api
func (ds *DefaultChecker) GetChecks(addr, path string) (Checks, error) {
	url := fmt.Sprintf("http://%s%s", addr, path)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	checks := Checks{}
	if err := json.Unmarshal(body, &checks); err != nil {
		return nil, err
	}
	return checks, nil
}

// GetChecks function to get a map of checks
func GetChecks(c Checker, addr string) (Checks, error) {
	checks, err := c.GetChecks(addr, checksPath)
	if err != nil {
		return nil, fmt.Errorf("Error querying Consul API: %s", err)
	}
	return checks, nil
}

// CheckRegister interface
type CheckRegister interface {
	Register(addr, path string, c *Check) error
	Unregister(addr, path, checkID string) error
}

// DefaultCheckRegister struct to implement CheckRegister default methods
type DefaultCheckRegister struct{}

// Register make PUT request on consul api
func (dr *DefaultCheckRegister) Register(addr, path string, c *Check) error {
	return register(addr, path, c)
}

// Unregister make PUT request on consul api
func (dr *DefaultCheckRegister) Unregister(addr, path, checkID string) error {
	return unregister(addr, path, checkID)
}

// DoServiceRegister function to register a new service
func DoCheckRegister(cr CheckRegister, addr string, c *Check) error {
	err := cr.Register(addr, registerCheckPath, c)
	if err != nil {
		return fmt.Errorf("Error registering check Consul API: %s", err)
	}
	return nil
}

// DoServiceUnregister function to unregister a service by ID
func DoCheckUnregister(cr CheckRegister, addr, checkID string) error {
	err := cr.Unregister(addr, deregisterServicePath, checkID)
	if err != nil {
		return fmt.Errorf("Error unregistering check Consul API: %s", err)
	}
	return nil
}

package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	agentServicesPath          = "/v1/agent/services"           // Returns the services the local agent is managing
	agentServiceRegisterPath   = "/v1/agent/service/register"   // Registers a new local service
	agentServiceDeregisterPath = "/v1/agent/service/deregister" // Deregisters a local service
)

// Service single consul service
type Service struct {
	ID      string   `json:"ID"`
	Service string   `json:"Service"`
	Name    string   `json:"Name"`
	Tags    []string `json:"Tags,omitempty"`
	Address string   `json:"Address,omitempty"`
	Port    int      `json:"Port"`
}

// Services map of services
type Services map[string]Service

// Servicer interface
type Servicer interface {
	GetServices(addr, path string) (Services, error)
}

// DefaultServicer struct to implement Servicer default methods
type DefaultServicer struct{}

// GetServices make GET request on consul api
func (ds *DefaultServicer) GetServices(addr, path string) (Services, error) {
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
	services := Services{}
	if err := json.Unmarshal(body, &services); err != nil {
		return nil, err
	}
	return services, nil
}

// GetServices function to get a map of services
func GetServices(s Servicer, addr string) (Services, error) {
	services, err := s.GetServices(addr, agentServicesPath)
	if err != nil {
		return nil, fmt.Errorf("Error querying Consul API: %s", err)
	}
	return services, nil
}

// ServiceRegister interface
type ServiceRegister interface {
	Register(addr, path string, s *Service) error
	Unregister(addr, path, serviceID string) error
}

// DefaultServiceRegister struct to implement Register default methods
type DefaultServiceRegister struct{}

// Register make PUT request on consul api
func (dr *DefaultServiceRegister) Register(addr, path string, s *Service) error {
	return register(addr, path, s)
}

// Unregister make PUT request on consul api
func (dr *DefaultServiceRegister) Unregister(addr, path, serviceID string) error {
	return unregister(addr, path, serviceID)
}

// DoServiceRegister function to register a new service
func DoServiceRegister(sr ServiceRegister, addr string, s *Service) error {
	err := sr.Register(addr, agentServiceRegisterPath, s)
	if err != nil {
		return fmt.Errorf("Error registering service Consul API: %s", err)
	}
	return nil
}

// DoServiceUnregister function to unregister a service by ID
func DoServiceUnregister(sr ServiceRegister, addr, serviceID string) error {
	err := sr.Unregister(addr, agentServiceDeregisterPath, serviceID)
	if err != nil {
		return fmt.Errorf("Error unregistering service Consul API: %s", err)
	}
	return nil
}

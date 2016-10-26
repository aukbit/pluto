package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	servicesPath          = "/v1/agent/services"                       // Returns the services the local agent is managing
	registerServicePath   = "/v1/agent/service/register"               // Registers a new local service
	deregisterServicePath = "/v1/agent/service/deregister/<serviceID>" // Deregisters a local service
)

// Service single consul service
type Service struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Tags    []string `json:"Tags,omitempty"`
	Address string   `json:"Address,omitempty"`
	Port    int      `json:"Port"`
}

// Services slice of service
type Services []Service

// Servicer interface
type Servicer interface {
	GetServices(addr, path string) (Services, error)
	Register(addr, path string)
	Deregister(addr, path string)
}

// DefaultServicer struct to implement Servicer default methods
type DefaultServicer struct{}

// GetServices make GET request on consul api
func (ds *DefaultServicer) GetServices(addr, path string) (Services, error) {
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
	services := Services{}
	if err := json.Unmarshal(body, &services); err != nil {
		return nil, err
	}
	return services, nil
}

// func services(url string) (map[string]*Service, error) {
//
// 	resp, err := http.Get(url + SERVICES)
// 	if err != nil {
// 		return nil, err
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data map[string]*Service
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
//
// func registerService(url string, s *Service) error {
//
// 	b, err := json.Marshal(s)
// 	if err != nil {
// 		return err
// 	}
// 	req, err := http.NewRequest("PUT", url+REGISTER_SERVICE, bytes.NewBuffer(b))
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return err
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("Error %v", string(body))
// 	}
// 	return nil
// }
//
// func deregisterService(url string, serviceID string) error {
//
// 	req, err := http.NewRequest("PUT", url+strings.Replace(DEREGISTER_SERVICE, "<serviceID>", serviceID, 1), bytes.NewBuffer([]byte(`{}`)))
// 	if err != nil {
// 		return err
// 	}
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	defer resp.Body.Close()
// 	if err != nil {
// 		return err
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("Error %v", string(body))
// 	}
// 	return nil
// }

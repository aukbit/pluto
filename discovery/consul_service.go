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
	Tags    []string `json:"Tags"`
	Address string   `json:"Address"`
	Port    int      `json:"Port"`
}

func Services() (map[string]*Service, error) {

	resp, err := http.Get(URL + SERVICES)
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

func RegisterService(s *Service) error {

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("PUT", URL+REGISTER_SERVICE, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// defer resp.Body.Close()
	// if err != nil {
	// 	return err
	// }
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error %v", resp.StatusCode)
	}
	return nil
}

func DeregisterService(serviceID string) error {

	req, err := http.NewRequest("PUT", URL+strings.Replace(DEREGISTER_SERVICE, "<serviceID>", serviceID, 1), bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// body, err := ioutil.ReadAll(resp.Body)
	// defer resp.Body.Close()
	// if err != nil {
	// 	log.Printf("RegisterService %v", err)
	// 	return err
	// }
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error %v", resp.StatusCode)
	}
	return nil
}

// func GetService(name string) (s *Service, err error) {
//
// 	qs := "?near=_agent"
// 	resp, err := http.Get(URL + strings.Replace(SERVICE, "<service>", name, 1) + qs)
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
// 		s := &Service{}
// 		if err := json.Unmarshal(bytes, s); err != nil {
// 			return nil, err
// 		}
// 		log.Printf("GetService %v", s)
// 		// returns the first it finds
// 		return s, nil
// 	}
// 	return &Service{}, nil
// }

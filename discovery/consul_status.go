package discovery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	statusLeaderPath = "/v1/status/leader" // Returns the current Raft leader
)

// Leader interface
type Leader interface {
	Status(addr, path string) (string, error)
}

// GetStatus function to get the status of consul leader node
func GetStatus(l Leader, addr string) (string, error) {
	status, err := l.Status(addr, statusLeaderPath)
	if err != nil {
		return "", fmt.Errorf("Error querying Consul API: %s", err)
	}
	return status, nil
}

// DefaultLeader struct to append Status
type DefaultLeader struct{}

// Status make GET request on consul api
func (dn *DefaultLeader) Status(addr, path string) (string, error) {
	url := fmt.Sprintf("http://%s%s", addr, path)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var target string
	if err := json.Unmarshal(body, &target); err != nil {
		return "", err
	}
	return target, nil
}

func isAvailable(addr string) (bool, error) {
	status, err := GetStatus(&DefaultLeader{}, addr)
	if status != "" {
		return true, nil
	}
	return false, err
}

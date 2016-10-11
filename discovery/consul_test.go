package discovery

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestNodes(t *testing.T) {

	nodes, err := Nodes()
	if err != nil {
		t.Error(err)
	}
	t.Logf("nodes %v", nodes)
}

func TestServices(t *testing.T) {

	services, err := Services()
	if err != nil {
		t.Error(err)
	}
	t.Logf("services %v", services)
}

func TestRegisterService(t *testing.T) {
	s := &Service{
		ID:   "test1",
		Name: "test1",
		Tags: []string{"auth", "api"},
		Port: 60500,
	}
	err := RegisterService(s)
	if err != nil {
		t.Error(err)
	}
	err = DeregisterService("test1")
	if err != nil {
		t.Error(err)
	}
}

func TestCatalogServices(t *testing.T) {
	services, err := CatalogServices()
	if err != nil {
		t.Error(err)
	}
	t.Logf("services %v", services)
}

func TestCatalogService(t *testing.T) {
	s := &Service{
		ID:   "test2",
		Name: "test2",
		Tags: []string{"auth", "api"},
		Port: 60500,
	}
	err := RegisterService(s)
	if err != nil {
		t.Error(err)
	}
	ns, err := CatalogService("test2")
	if err != nil {
		t.Error(err)
	}
	t.Logf("services %v", ns)
	err = DeregisterService("test2")
	if err != nil {
		t.Error(err)
	}
}

func TestIsAvailable(t *testing.T) {
	ok, err := IsAvailable()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, ok)
}

func TestRegisterCheck(t *testing.T) {
	s := &Service{
		ID:   "test3",
		Name: "test3",
		Tags: []string{"auth", "api"},
		Port: 60500,
	}
	err := RegisterService(s)
	if err != nil {
		t.Error(err)
	}
	c := &Check{
		ID:    "test3_check",
		Name:  "TCP check",
		Notes: "Ensure the server is listening on the specific port",
		DeregisterCriticalServiceAfter: "1m",
		TCP:       ":60500",
		Interval:  "10s",
		Timeout:   "1s",
		ServiceID: "test3",
	}
	err = RegisterCheck(c)
	if err != nil {
		t.Error(err)
	}
}

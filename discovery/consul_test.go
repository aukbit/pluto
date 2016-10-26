package discovery

import (
	"testing"

	"github.com/paulormart/assert"
)

const URL = "http://192.168.99.100:8500"

func TestServices(t *testing.T) {

	services, err := services(URL)
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
	err := registerService(URL, s)
	if err != nil {
		t.Error(err)
	}
	err = deregisterService(URL, "test1")
	if err != nil {
		t.Error(err)
	}
}

func TestCatalogServices(t *testing.T) {
	services, err := CatalogServices(URL)
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
	err := registerService(URL, s)
	if err != nil {
		t.Error(err)
	}
	ns, err := CatalogService(URL, "test2")
	if err != nil {
		t.Error(err)
	}
	t.Logf("services %v", ns)
	err = deregisterService(URL, "test2")
	if err != nil {
		t.Error(err)
	}
}

func TestIsAvailable(t *testing.T) {
	ok, err := isAvailable(URL)
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
	err := registerService(URL, s)
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
	err = registerCheck(URL, c)
	if err != nil {
		t.Error(err)
	}
}

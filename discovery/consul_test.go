package discovery

import "testing"

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
		ID:   "test",
		Name: "test",
		Tags: []string{"auth", "api"},
		Port: 60500,
	}
	err := RegisterService(s)
	if err != nil {
		t.Error(err)
	}
}

func TestDeregisterService(t *testing.T) {
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
		ID:   "test",
		Name: "test",
		Tags: []string{"auth", "api"},
		Port: 60500,
	}
	err := RegisterService(s)
	if err != nil {
		t.Error(err)
	}
	ns, err := CatalogService("test")
	if err != nil {
		t.Error(err)
	}
	t.Logf("services %v", ns)
	err = DeregisterService("test")
	if err != nil {
		t.Error(err)
	}
}

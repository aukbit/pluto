package discovery

import (
	"testing"

	"github.com/paulormart/assert"
)

const URL = "http://192.168.99.100:8500"

// func _TestServices(t *testing.T) {
//
// 	services, err := services(URL)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Logf("services %v", services)
// }

// func _TestCatalogServices(t *testing.T) {
// 	services, err := CatalogServices(URL)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Logf("services %v", services)
// }

// func _TestCatalogService(t *testing.T) {
// 	s := &Service{
// 		ID:      "test2",
// 		Service: "test2",
// 		Tags:    []string{"auth", "api"},
// 		Port:    60500,
// 	}
// 	err := registerService(URL, s)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	ns, err := CatalogService(URL, "test2")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Logf("services %v", ns)
// 	err = deregisterService(URL, "test2")
// 	if err != nil {
// 		t.Error(err)
// 	}
// }

func _TestIsAvailable(t *testing.T) {
	ok, err := isAvailable(URL)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, true, ok)
}

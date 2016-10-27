package discovery

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestDefaults(t *testing.T) {
	c := newConfig()
	assert.Equal(t, "localhost:8500", c.Addr)
}

func TestConfigs(t *testing.T) {
	s := &Service{
		ID:      "test",
		Service: "test",
		Tags:    []string{"auth", "api"},
		Port:    60500,
	}
	c := &Check{
		ID:    "test_check",
		Name:  "TCP check",
		Notes: "Ensure the server is listening on the specific port",
		DeregisterCriticalServiceAfter: "1m",
		TCP:       ":60500",
		Interval:  "10s",
		Timeout:   "1s",
		ServiceID: "test",
	}
	cfg := newConfig(
		Addr("192.168.99.100:8500"),
		ServicesCfg(s),
		ChecksCfg(c),
	)
	assert.Equal(t, "192.168.99.100:8500", cfg.Addr)
	assert.Equal(t, "http://192.168.99.100:8500", cfg.URL())
	assert.Equal(t, []*Service{s}, cfg.Services)
	assert.Equal(t, []*Check{c}, cfg.Checks)
}

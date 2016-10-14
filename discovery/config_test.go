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
	c := newConfig(Addr("192.168.99.100:8500"))
	assert.Equal(t, "192.168.99.100:8500", c.Addr)
}

package client

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestDefaults(t *testing.T) {
	c := newConfig()
	assert.Equal(t, true, len(c.ID) > 0)
	assert.Equal(t, "", c.Description)
	assert.Equal(t, "client", c.Name)
	assert.Equal(t, "grpc", c.Format)
	assert.Equal(t, "1.0.0", c.Version)
	assert.Equal(t, "localhost:65060", c.Target)
}

func TestConfigs(t *testing.T) {
	c := newConfig(ID("123456"), Name("Special"), Target("localhost:65062"))
	assert.Equal(t, "123456", c.ID)
	// Note: lower case and prefix 'server_' in name
	assert.Equal(t, "client_special", c.Name)
	assert.Equal(t, "localhost:65062", c.Target)
}

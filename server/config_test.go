package server

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestDefaults(t *testing.T) {
	c := newConfig()
	assert.Equal(t, true, len(c.ID) > 0)
	assert.Equal(t, ":8080", c.Addr)
	assert.Equal(t, "", c.Description)
	assert.Equal(t, "server", c.Name)
	assert.Equal(t, "http", c.Format)
	assert.Equal(t, "v1.0.0", c.Version)
	assert.Equal(t, nil, c.Mux)
}

func TestConfigs(t *testing.T) {
	c := newConfig(
		ID("123456"),
		Name("Special"),
		Description("Special server description"),
		Addr(":8081"))
	assert.Equal(t, "123456", c.ID)
	// Note: lower case and prefix 'server_' in name
	assert.Equal(t, "server_special", c.Name)
	assert.Equal(t, "Special server description", c.Description)
	assert.Equal(t, ":8081", c.Addr)
}

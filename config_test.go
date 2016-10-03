package pluto

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestDefaults(t *testing.T) {
	c := newConfig()
	assert.Equal(t, true, len(c.ID) > 0)
	assert.Equal(t, "", c.Description)
	assert.Equal(t, "pluto", c.Name)
	assert.Equal(t, "1.0.0", c.Version)
}

func TestConfigs(t *testing.T) {
	c := newConfig(ID("123456"), Name("Special"))
	assert.Equal(t, "123456", c.ID)
	// Note: lower case and prefix 'server_' in name
	assert.Equal(t, "pluto_special", c.Name)
}

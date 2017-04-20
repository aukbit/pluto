package pluto

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestConfig(t *testing.T) {
	c := newConfig()
	assert.Equal(t, true, len(c.ID) > 0)
	assert.Equal(t, "", c.Description)
	assert.Equal(t, "pluto", c.Name)
	assert.Equal(t, "1.3.2", c.Version)
}

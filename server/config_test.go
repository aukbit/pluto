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
	assert.Equal(t, "1.3.1", c.Version)
	assert.Equal(t, nil, c.Mux)
}

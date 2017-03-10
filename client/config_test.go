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
	assert.Equal(t, "v1.0.0", c.Version)
	assert.Equal(t, "localhost:65060", c.Target())
}

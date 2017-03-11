package datastore

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestDefaults(t *testing.T) {
	c := newConfig()
	assert.Equal(t, "default", c.Keyspace)
	assert.Equal(t, "127.0.0.1", c.Target)
}

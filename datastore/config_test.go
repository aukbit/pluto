package datastore

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestDefaults(t *testing.T) {
	c := newConfig()
	assert.Equal(t, "default", c.Keyspace)
	assert.Equal(t, "127.0.0.1", c.Addr)
}

func TestConfigs(t *testing.T) {
	c := newConfig(Keyspace("mars"), Addr("localhost"))
	assert.Equal(t, "mars", c.Keyspace)
	assert.Equal(t, "localhost", c.Addr)
}

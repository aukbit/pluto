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

func TestConfigs(t *testing.T) {
	c := newConfig(Keyspace("mars"), Target("localhost"))
	assert.Equal(t, "mars", c.Keyspace)
	assert.Equal(t, "localhost", c.Target)
}

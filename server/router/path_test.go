package router

import (
	"testing"

	"github.com/paulormart/assert"
)

func TestPaths(t *testing.T) {

	var path, key, value, prefix string
	path = "/"
	key, value, prefix, params := transformPath(path)
	assert.Equal(t, "/", key)
	assert.Equal(t, "/", value)
	assert.Equal(t, "", prefix)
	assert.Equal(t, 0, len(params))
	path = "/home"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home", key)
	assert.Equal(t, "/home", value)
	assert.Equal(t, "", prefix)
	assert.Equal(t, 0, len(params))
	path = "/home/:id"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home/:", key)
	assert.Equal(t, "/:", value)
	assert.Equal(t, "/home", prefix)
	assert.Equal(t, "id", params[0])
	path = "/home/:id/room"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home/:/room", key)
	assert.Equal(t, "/room", value)
	assert.Equal(t, "/home/:", prefix)
	assert.Equal(t, "id", params[0])

}

package router

import (
	"assert"
	"testing"
)

func TestPaths(t *testing.T) {

	var path, key, value, prefix string
	var params []string
	path = "/"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/", key)
	assert.Equal(t, "/", value)
	assert.Equal(t, "", prefix)
	assert.Equal(t, []string{}, params)
	path = "/home"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home", key)
	assert.Equal(t, "/home", value)
	assert.Equal(t, "", prefix)
	assert.Equal(t, []string{}, params)
	path = "/home/:id"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home/:", key)
	assert.Equal(t, "/:", value)
	assert.Equal(t, "/home", prefix)
	assert.Equal(t, []string{"id"}, params)
	path = "/home/:id/room"
	key, value, prefix, params = transformPath(path)
	assert.Equal(t, "/home/:/room", key)
	assert.Equal(t, "/room", value)
	assert.Equal(t, "/home/:", prefix)
	assert.Equal(t, []string{"id"}, params)

}

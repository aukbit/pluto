package client_test

import (
	"testing"
	"pluto/client"
	"github.com/paulormart/assert"
	"reflect"
)

func TestClient(t *testing.T){

	//1. create a client
	c := client.NewClient(
		client.Name("gopher"),
		client.Description("gopher super client"),
	)
	assert.Equal(t, reflect.TypeOf(client.DefaultClient), reflect.TypeOf(c))

	cfg := c.Config()
	assert.Equal(t, true, len(cfg.Id) > 0)
	assert.Equal(t, "gopher.client", cfg.Name)
	assert.Equal(t, "gopher super client", cfg.Description)

}
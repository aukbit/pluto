package client_test

import (
	"testing"
)

func TestClient(t *testing.T){

	//1. create a client
	c := client.NewClient(
		client.Name("gopher"),
		client.Description("gopher super server"),
		client.Addr(":8080"),
	)
}
package pluto

import (
	context "golang.org/x/net/context"

	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/discovery"
	"github.com/aukbit/pluto/server"
)

// Config pluto service config
type Config struct {
	ID          string
	Name        string
	Description string
	Version     string
	Servers     map[string]*server.Server
	Clients     map[string]*client.Client
	clientsCh   chan *client.Client
	Datastore   *datastore.Datastore
	Discovery   discovery.Discovery
	Hooks       map[string][]HookFunc
	HealthAddr  string // TCP address (e.g. localhost:8000) to listen on, ":http" if empty
}

// HookFunc hook function type
type HookFunc func(context.Context)

func newConfig() Config {
	return Config{
		ID:         common.RandID("plt_", 6),
		Name:       defaultName,
		Version:    defaultVersion,
		Servers:    make(map[string]*server.Server),
		Clients:    make(map[string]*client.Client),
		clientsCh:  make(chan *client.Client, 100),
		Hooks:      make(map[string][]HookFunc),
		HealthAddr: defaultHealthAddr,
	}
}

package pluto

import (
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/grpc/health"

	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/datastore"
	"bitbucket.org/aukbit/pluto/server"
	"github.com/uber-go/zap"
)

var logger = zap.New(zap.NewJSONEncoder())

// Config pluto service config
type Config struct {
	ID          string
	Name        string
	Description string
	Version     string
	Servers     map[string]server.Server
	Clients     map[string]client.Client
	Datastore   datastore.Datastore
	health      *health.Server
}

// ConfigFunc registers the Config
type ConfigFunc func(*Config)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Version: defaultVersion,
		Servers: make(map[string]server.Server),
		Clients: make(map[string]client.Client),
		health:  health.NewServer()}

	for _, c := range cfgs {
		c(cfg)
	}

	if len(cfg.ID) == 0 {
		cfg.ID = common.RandID("plt_", 6)
	}

	if len(cfg.Name) == 0 {
		cfg.Name = defaultName
	}

	return cfg
}

// ID service id
func ID(id string) ConfigFunc {
	return func(cfg *Config) {
		cfg.ID = id
	}
}

// Name service name
func Name(n string) ConfigFunc {
	return func(cfg *Config) {
		// support only alphanumeric and underscore characters
		reg, err := regexp.Compile("[^A-Za-z0-9_]+")
		if err != nil {
			logger.Error("Name",
				zap.String("err", err.Error()),
			)
		}
		safe := reg.ReplaceAllString(n, "_")
		cfg.Name = fmt.Sprintf("%s_%s", defaultName, strings.ToLower(safe))
	}
}

// Description service description
func Description(d string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Description = d
	}
}

// Servers slice of service servers
func Servers(s server.Server) ConfigFunc {
	return func(cfg *Config) {
		cfg.Servers[s.Config().Name] = s
		cfg.health.SetServingStatus(s.Config().Name, 2)
	}
}

// Clients slice of service clients
func Clients(c client.Client) ConfigFunc {
	return func(cfg *Config) {
		cfg.Clients[c.Config().Name] = c
		cfg.health.SetServingStatus(c.Config().Name, 2)
	}
}

// Datastore to persist data
func Datastore(addr string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Datastore = datastore.NewDatastore(datastore.Addr(addr), datastore.Keyspace(cfg.Name))
	}
}

// Datastore to persist data
func DatastoreDiscovery(service string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Datastore = datastore.NewDatastore(datastore.TargetDiscovery(service), datastore.Keyspace(cfg.Name))
	}
}

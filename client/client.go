package client

import "github.com/uber-go/zap"

// Client is an interface to make calls to services
type Client interface {
	Dial() error
	Call() interface{}
	Close() error
	Config() *Config
}

var (
	defaultName    = "client"
	defaultVersion = "1.0.0"
	logger         = zap.New(zap.NewJSONEncoder())
)

// NewClient returns a new client with cfg passed in
func NewClient(cfgs ...ConfigFunc) Client {
	return newClient(cfgs...)
}

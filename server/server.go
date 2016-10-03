package server

import "github.com/uber-go/zap"

// Server is the basic interface that defines what to expect from any server.
type Server interface {
	//Init(...ConfigFunc)		error	// TODO remove init! tere is no need
	Run() error
	Stop()
	Config() *Config
}

var (
	defaultName    = "server"
	defaultVersion = "1.0.0"
	logger         = zap.New(zap.NewJSONEncoder())
)

// NewServer returns a new http server with cfg passed in
func NewServer(cfgs ...ConfigFunc) Server {
	return newServer(cfgs...)
}

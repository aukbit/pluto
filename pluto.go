package pluto

import (
	"pluto/server"
	"pluto/client"
)


// Service is the basic interface that defines what to expect from any server.
type Service interface {
	Init(...ConfigFunc) 		error
	Servers() 			map[string]server.Server
	Clients() 			map[string]client.Client
	Run() 				error
	Stop() 				error
	Config() 			*Config
}

var (
	DefaultName 		= "pluto"
	DefaultVersion      	= "1.0.0"
)

// NewService returns a new service with cfg passed in
func NewService(cfgs ...ConfigFunc) Service {
	return newService(cfgs...)
}
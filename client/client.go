package client

import (
	"github.com/google/uuid"
)

// Client is an interface to make calls to services
type Client interface {
	Init(...ConfigFunc)						error
	Run() 									error
	Stop() 									error
	Config() 								*Config
}

var (
	DefaultName			= "client"
	DefaultVersion      = "1.0.0"
	DefaultId			= uuid.New().String()
	DefaultClient  		= newDefaultClient()
)

// NewClient returns a new client with cfg passed in
func NewClient(cfgs ...ConfigFunc) Client {
	return newDefaultClient(cfgs...)
}
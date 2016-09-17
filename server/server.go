package server

import (
	"log"
)

// Server is the basic interface that defines what to expect from any server.
type Server interface {
	Init(...ConfigFunc)		error
	Run() 				error
	Stop() 				error
	Config() 			*Config
}

var (
	DefaultName			= "server"
	DefaultVersion      		= "1.0.0"
	DefaultServer  			= newDefaultServer()
)

// NewServer returns a new http server with cfg passed in
func NewServer(cfgs ...ConfigFunc) Server {
	return newDefaultServer(cfgs...)
}

// NewGRPCServer returns a new grpc server with cfg passed in
func NewGRPCServer(cfgs ...ConfigFunc) Server {
	return newGRPCServer(cfgs...)
}

// Init initialises the default server with options passed in
func Init(cfgs ...ConfigFunc) {
	if DefaultServer == nil {
		DefaultServer = newDefaultServer(cfgs...)
	}
	DefaultServer.Init(cfgs...)
}

// Run will start a DefaultServer and set it up to Stop()
// on a kill signal.
func Run() error {
	log.Printf("Run server")
	if err := DefaultServer.Run(); err != nil {
		return err
	}
	return nil
}

// Stop stops the default server
func Stop() error {
	log.Printf("Stop server")
	if err := DefaultServer.Stop(); err != nil {
		return err
	}
	return nil
}

package server

import (
	"log"
	"net"

)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type grpcServer struct {
	defaultServer
}

// serve serves *grpc.Server
func (s *grpcServer) serve(ln net.Listener) (err error) {

	srv := s.cfg.GRPCServer

	go func() {
		if err := srv.Serve(ln); err != nil {
			log.Fatalf("ERROR %s g.Serve(lis) %v", s.cfg.Name, err)
		}
	}()

	log.Printf("----- %s %s listening on %s", s.cfg.Format, s.cfg.Name, ln.Addr().String())
	return nil
}

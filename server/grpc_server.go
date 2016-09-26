package server

import (
	"log"
	"net"
	"syscall"
	"os/signal"
	"os"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type grpcServer struct {
	defaultServer
}

// Run Server
func (s *grpcServer) Run() error {
	if err := s.start(); err != nil {
		return err
	}
	// parse address for host, port
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	sig := <-ch
	log.Printf("----- %s signal %v received ", s.cfg.Name, sig)
	return s.Stop()
}

func (s *grpcServer) start() error {

	ln, err := s.listen()
	if err != nil {
		return err
	}

	if err := s.serve(ln); err != nil {
		return err
	}
	go s.waitSignal(ln)
	return nil
}

// listen based on http.ListenAndServe
// listens on the TCP network address srv.Addr
// If srv.Addr is blank, ":http" is used.
// returns nil or new listener
func (s *grpcServer) listen() (net.Listener, error) {

	addr := s.cfg.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	ln = net.Listener(TcpKeepAliveListener{ln.(*net.TCPListener)})

	return ln, nil
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

package server

import (
	"net"

	"google.golang.org/grpc"

	"github.com/uber-go/zap"
)

type B struct{}

// serveGRPC serves *grpc.Server
func (ds *defaultServer) serveGRPC(ln net.Listener) (err error) {

	// register
	//
	srv := grpc.NewServer(grpc.UnaryInterceptor(WrapUnaryInterceptor("logger", ds.logger)))
	// srv := grpc.NewServer()
	ds.cfg.GRPCRegister(srv)

	// srv := ds.cfg.GRPCServer
	// add go routine to WaitGroup
	ds.wg.Add(1)
	go func() {
		defer ds.wg.Done()
		if err := srv.Serve(ln); err != nil {
			if err.Error() == errClosing(ln).Error() {
				return
			}
			ds.logger.Error("Serve(ln)", zap.String("err", err.Error()))
			return
		}
	}()
	return nil
}

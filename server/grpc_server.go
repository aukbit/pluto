package server

import (
	"net"

	"google.golang.org/grpc"

	"github.com/uber-go/zap"
)

// serveGRPC serves *grpc.Server
func (ds *defaultServer) serveGRPC(ln net.Listener) (err error) {

	// add go routine to WaitGroup
	ds.wg.Add(1)
	go func(srv *grpc.Server) {
		defer ds.wg.Done()
		if err := srv.Serve(ln); err != nil {
			if err.Error() == errClosing(ln).Error() {
				return
			}
			ds.logger.Error("Serve(ln)", zap.String("err", err.Error()))
			return
		}
	}(ds.grpcServer)
	return nil
}

package server

import (
	"net"

	"github.com/uber-go/zap"
)

// serveGRPC serves *grpc.Server
func (s *defaultServer) serveGRPC(ln net.Listener) (err error) {

	srv := s.cfg.GRPCServer
	// add go routine to WaitGroup
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := srv.Serve(ln); err != nil {
			if err.Error() == errClosing(ln).Error() {
				return
			}
			logger.Error("Serve(ln)",
				zap.String("server", s.cfg.Name),
				zap.String("id", s.cfg.ID),
				zap.String("format", s.cfg.Format),
				zap.String("port", ln.Addr().String()),
				zap.String("err", err.Error()))
			return
		}
	}()
	return nil
}

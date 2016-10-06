package server

import (
	"bitbucket.org/aukbit/pluto/common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// getOrCreateEventID uses grpc metadata context to set an event id
// the metadata context is then sent over the wire - gRPC calls
// and available to other services
func getOrCreateEventID(ctx context.Context) (string, context.Context) {
	// get
	md, ok := metadata.FromContext(ctx)
	if ok {
		e, ok := md["event"]
		if ok {
			return e[0], ctx
		}
	}
	// create
	e := common.RandID("evt_", 12)
	ctx = metadata.NewContext(ctx, metadata.Pairs("event", e))
	return e, ctx
}

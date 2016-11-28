package common

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// GetOrCreateEventID uses grpc metadata context to set an event id
// the metadata context is then sent over the wire - gRPC calls
// and available to other services
func GetOrCreateEventID(ctx context.Context) (string, context.Context) {
	// get metadata from context
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	e, ok := md["event"]
	if !ok {
		// append new evt id
		md["event"] = append(md["event"], RandID("evt_", 12))
	}
	ctx = metadata.NewContext(ctx, md)
	e, _ = md["event"]
	// log.Printf("GetOrCreateEventID md:%v", md)
	// log.Printf("GetOrCreateEventID ctx:%v", ctx)

	return e[0], ctx
}

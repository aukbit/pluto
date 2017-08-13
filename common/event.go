package common

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

// GetOrCreateEventID uses grpc metadata context to set an event id
// the metadata context is then sent over the wire - gRPC calls
// and available to other services
func DEPRECATEDGetOrCreateEventID(ctx context.Context) (string, context.Context) {
	// get metadata from context
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	e, ok := md["event"]
	if !ok {
		// append new evt id
		md["event"] = append(md["event"], RandID("evt_", 16))
	}
	ctx = metadata.NewContext(ctx, md)
	e, _ = md["event"]
	// log.Printf("GetOrCreateEventID md:%v", md)
	// log.Printf("GetOrCreateEventID ctx:%v", ctx)

	return e[0], ctx
}

// GetOrCreateEventID returns eid and context with eid in context metadata
func GetOrCreateEventID(ctx context.Context) (string, context.Context) {
	// get metadata from context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	// NOTE: modification should be made to copies of the returned MD.
	md = md.Copy()
	// verify if eid exists, if not generate new eid
	_, ok = md["eid"]
	if !ok {
		md = metadata.Join(md, metadata.Pairs("eid", RandID("", 16)))
	}
	ctx = metadata.NewOutgoingContext(ctx, md)
	return md["eid"][0], ctx
}

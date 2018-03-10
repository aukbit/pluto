package common

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

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

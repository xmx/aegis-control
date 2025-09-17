package linkhub

import (
	"context"

	"github.com/xmx/aegis-common/transport"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func WithValue(parent context.Context, p transport.Peer[bson.ObjectID]) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	if p == nil {
		return parent
	}

	return context.WithValue(parent, defaultContextKey, p)
}

func FromContext(ctx context.Context) transport.Peer[bson.ObjectID] {
	if ctx == nil {
		return nil
	}

	val := ctx.Value(defaultContextKey)
	if p, ok := val.(transport.Peer[bson.ObjectID]); ok {
		return p
	}

	return nil
}

var defaultContextKey = contextKey{}

type contextKey struct{}

func (contextKey) String() string {
	return "transport-peer-context-key"
}

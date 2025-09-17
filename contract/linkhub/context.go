package linkhub

import "context"

func WithValue(parent context.Context, p Peer) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	if p == nil {
		return parent
	}

	return context.WithValue(parent, defaultContextKey, p)
}

func FromContext(ctx context.Context) Peer {
	if ctx == nil {
		return nil
	}

	val := ctx.Value(defaultContextKey)
	if p, ok := val.(Peer); ok {
		return p
	}

	return nil
}

var defaultContextKey = contextKey{}

type contextKey struct{}

func (contextKey) String() string {
	return "transport-peer-context-key"
}

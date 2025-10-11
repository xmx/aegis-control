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

func FromContext(ctx context.Context) (Peer, bool) {
	if ctx == nil {
		return nil, false
	}

	val := ctx.Value(defaultContextKey)
	if p, _ := val.(Peer); p != nil {
		return p, true
	}

	return nil, false
}

var defaultContextKey = contextKey{}

type contextKey struct{}

func (contextKey) String() string {
	return "tunnel-peer-context-key"
}

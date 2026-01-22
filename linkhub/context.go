package linkhub

import "context"

func WithValue(parent context.Context, p Peer) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	if p == nil {
		return parent
	}

	return context.WithValue(parent, peerContextKey, p)
}

func FromContext(ctx context.Context) (Peer, bool) {
	if ctx == nil {
		return nil, false
	}

	val := ctx.Value(peerContextKey)
	if p, _ := val.(Peer); p != nil {
		return p, true
	}

	return nil, false
}

type contextKey struct {
	name string
}

var peerContextKey = &contextKey{name: "tunnel-peer"}

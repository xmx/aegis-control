package linkhub

import "context"

func WithValue(parent context.Context, p Peer) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	if p == nil {
		return parent
	}

	return context.WithValue(parent, httpContextKey, p)
}

func FromContext(ctx context.Context) (Peer, bool) {
	if ctx == nil {
		return nil, false
	}

	val := ctx.Value(httpContextKey)
	if p, _ := val.(Peer); p != nil {
		return p, true
	}

	return nil, false
}

var httpContextKey = tunnelContextKey{name: "http-context-key"}

type tunnelContextKey struct {
	name string
}

func (k tunnelContextKey) String() string { return k.name }

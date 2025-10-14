package linkhub

import (
	"context"
	"net"
	"strings"

	"github.com/xmx/aegis-common/tunnel/tunutil"
)

func NewSuffixDialer(hub Huber, suffix string) tunutil.DialMatcher {
	return &suffixDialer{
		suffix: suffix,
		huber:  hub,
	}
}

type suffixDialer struct {
	suffix string
	huber  Huber
}

func (sd *suffixDialer) MatchDialer(_, address string) tunutil.Dialer {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil
	}

	if strings.HasSuffix(host, sd.suffix) {
		return sd
	}

	return nil
}

func (sd *suffixDialer) DialContext(ctx context.Context, _, address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	id, found := strings.CutSuffix(host, sd.suffix)
	if !found {
		return nil, &net.AddrError{
			Addr: address,
		}
	}

	mux := sd.huber.Get(id)
	if mux == nil {
		return nil, &net.AddrError{
			Addr: address,
		}
	}

	return mux.Muxer().Open(ctx)
}

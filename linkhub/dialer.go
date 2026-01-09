package linkhub

import (
	"context"
	"net"
	"strings"

	"github.com/xmx/aegis-common/muxlink/muxproto"
)

func NewSuffixDialer(suffix string, hub Huber) muxproto.Dialer {
	return &suffixDialer{
		huber:  hub,
		suffix: suffix,
	}
}

type suffixDialer struct {
	suffix string
	huber  Huber
}

func (sd *suffixDialer) Dial(network, address string) (net.Conn, error) {
	return sd.DialContext(context.Background(), network, address)
}

func (sd *suffixDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, _, _ := net.SplitHostPort(address)
	if host == "" {
		return nil, nil
	}
	id, found := strings.CutSuffix(host, sd.suffix)
	if !found {
		return nil, nil
	}

	if peer := sd.huber.Get(id); peer != nil {
		mux := peer.Muxer()
		return mux.Open(ctx)
	}

	return nil, &net.OpError{
		Op:   "dial",
		Net:  network,
		Addr: &net.UnixAddr{Net: network, Name: address},
		Err:  net.UnknownNetworkError("no route to host"),
	}
}

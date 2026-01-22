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

func NewMixedDialer(mux muxproto.MUXOpener, hub Huber, back muxproto.Dialer) muxproto.Dialer {
	return &mixedDialer{
		mux:  mux,
		hub:  hub,
		back: back,
	}
}

type mixedDialer struct {
	mux  muxproto.MUXOpener
	hub  Huber
	back muxproto.Dialer
}

func (m *mixedDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	if m.mux != nil && host == m.mux.Host() {
		return m.mux.Open(ctx)
	}

	if m.hub != nil {
		_, domain, found := strings.Cut(host, ".")
		if found && domain == m.hub.Domain() {
			peer := m.hub.Get(host)
			if peer == nil {
				return nil, peerUnreachable(network, host)
			}
			mux := peer.Muxer()

			return mux.Open(ctx)
		}
	}

	if m.back != nil {
		return m.back.DialContext(ctx, network, address)
	}

	return nil, &net.OpError{
		Op:   "dial",
		Net:  network,
		Addr: &net.UnixAddr{Net: network, Name: address},
		Err:  net.UnknownNetworkError("没有找到任何拨号器"),
	}
}

func peerUnreachable(network, address string) error {
	return &net.OpError{
		Op:   "lookup",
		Net:  "tunnel",
		Addr: &net.UnixAddr{Net: network, Name: address},
		Err:  net.UnknownNetworkError("节点未上线或未注册"),
	}
}

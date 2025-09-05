package quick

import (
	"context"
	"net"
	"time"

	"github.com/xmx/aegis-common/library/timex"
	"github.com/xmx/aegis-common/transport"
	"golang.org/x/net/quic"
)

type Server struct {
	Addr       string
	Handler    transport.Handler
	QUICConfig *quic.Config
	endpoint   *quic.Endpoint
}

func (s *Server) Serve(ctx context.Context, endpoint *quic.Endpoint) error {
	s.endpoint = endpoint

	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err := endpoint.Accept(ctx)
		if err != nil {
			if exx := ctx.Err(); exx != nil {
				return exx
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				tempDelay = min(tempDelay, time.Second)
				_ = timex.Sleep(ctx, tempDelay)
				continue
			}
		}

		go s.serve(ctx, conn)
	}
}

func (s *Server) Close() error {
	return s.endpoint.Close(context.Background())
}

func (s *Server) serve(parent context.Context, conn *quic.Conn) {
	mux := transport.NewQUIC(parent, conn, nil)
	_ = s.Handler.Handle(mux)
}

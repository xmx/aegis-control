package quick

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/xmx/aegis-common/muxlink/muxconn"
	"github.com/xmx/aegis-common/muxlink/muxproto"
)

type Server interface {
	io.Closer
	ListenAndServe(ctx context.Context) error
}

type QUICgo struct {
	Addr       string
	Handler    muxproto.MUXAccepter
	TLSConfig  *tls.Config
	QUICConfig *quic.Config
	mutex      sync.Mutex
	listeners  map[*quic.Listener]struct{}
}

func (q *QUICgo) Close() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	errs := make([]error, 0, len(q.listeners))
	for lis := range q.listeners {
		err := lis.Close()
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (q *QUICgo) ListenAndServe(ctx context.Context) error {
	addr := q.Addr
	if addr == "" {
		addr = ":443"
	}

	lis, err := quic.ListenAddr(addr, q.TLSConfig, q.QUICConfig)
	if err != nil {
		return err
	}
	defer lis.Close()

	q.mutex.Lock()
	if q.listeners == nil {
		q.listeners = make(map[*quic.Listener]struct{}, 4)
	}
	q.listeners[lis] = struct{}{}
	q.mutex.Unlock()

	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err1 := lis.Accept(ctx)
		if err1 != nil {
			if exx := ctx.Err(); exx != nil {
				return exx
			}
			if ne, ok := err1.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				tempDelay = min(tempDelay, time.Second)
				_ = timeSleep(ctx, tempDelay)
				continue
			}
		}

		go q.handle(ctx, conn)
	}
}

func (q *QUICgo) handle(parent context.Context, conn *quic.Conn) {
	mux := muxconn.NewQUICgo(parent, conn)
	if h := q.Handler; h != nil {
		h.AcceptMUX(mux)
	}
	_ = mux.Close()
}

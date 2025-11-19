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
	"github.com/xmx/aegis-common/library/timex"
	"github.com/xmx/aegis-common/tunnel/tundial"
	"github.com/xmx/aegis-common/tunnel/tunutil"
)

type Server interface {
	io.Closer
	ListenAndServe(ctx context.Context) error
}

type QUICGo struct {
	Addr       string
	Handler    tunutil.Handler
	TLSConfig  *tls.Config
	QUICConfig *quic.Config
	mutex      sync.Mutex
	listeners  map[*quic.Listener]struct{}
}

func (q *QUICGo) Close() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	errs := make([]error, 0, len(q.listeners))
	for lis := range q.listeners {
		err := lis.Close()
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (q *QUICGo) ListenAndServe(ctx context.Context) error {
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
				_ = timex.Sleep(ctx, tempDelay)
				continue
			}
		}

		go q.handle(ctx, conn)
	}
}

func (q *QUICGo) handle(parent context.Context, conn *quic.Conn) {
	mux := tundial.NewQUICGo(parent, conn)
	if h := q.Handler; h != nil {
		h.Handle(mux)
	}
	_ = mux.Close()
}

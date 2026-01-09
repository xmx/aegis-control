package quick

import (
	"context"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/xmx/aegis-common/library/timex"
	"github.com/xmx/aegis-common/muxlink/muxconn"
	"github.com/xmx/aegis-common/muxlink/muxproto"
	"golang.org/x/net/quic"
)

type QUICx struct {
	Addr       string
	Accept     muxproto.MUXAccepter
	QUICConfig *quic.Config
	mutex      sync.Mutex
	endpoints  map[*quic.Endpoint]struct{}
}

func (q *QUICx) Close() error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	errs := make([]error, 0, len(q.endpoints))
	for end := range q.endpoints {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := end.Close(ctx)
		cancel()
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (q *QUICx) ListenAndServe(ctx context.Context) error {
	addr := q.Addr
	if addr == "" {
		addr = ":443"
	}

	endpoint, err := quic.Listen("udp", addr, q.QUICConfig)
	if err != nil {
		return err
	}
	q.mutex.Lock()
	if q.endpoints == nil {
		q.endpoints = make(map[*quic.Endpoint]struct{}, 4)
	}
	q.endpoints[endpoint] = struct{}{}
	q.mutex.Unlock()

	defer func() {
		ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
		_ = endpoint.Close(ctx1)
		cancel()
	}()

	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		conn, err1 := endpoint.Accept(ctx)
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

func (q *QUICx) handle(parent context.Context, conn *quic.Conn) {
	mux := muxconn.NewQUICx(parent, nil, conn)
	if acpt := q.Accept; acpt != nil {
		acpt.AcceptMUX(mux)
	}
	_ = mux.Close()
}

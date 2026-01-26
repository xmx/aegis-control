package telemetry

import (
	"context"
	"sync/atomic"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type spanProcessor struct {
	ptr atomic.Pointer[sdktrace.SpanProcessor]
}

func (sp *spanProcessor) OnStart(parent context.Context, s sdktrace.ReadWriteSpan) {
	if p := sp.load(); p != nil {
		p.OnStart(parent, s)
	}
}

func (sp *spanProcessor) OnEnd(s sdktrace.ReadOnlySpan) {
	if p := sp.load(); p != nil {
		p.OnEnd(s)
	}
}

func (sp *spanProcessor) Shutdown(ctx context.Context) error {
	if p := sp.load(); p != nil {
		return p.Shutdown(ctx)
	}

	return nil
}

func (sp *spanProcessor) ForceFlush(ctx context.Context) error {
	if p := sp.load(); p != nil {
		return p.ForceFlush(ctx)
	}

	return nil
}

func (sp *spanProcessor) swap(p sdktrace.SpanProcessor) sdktrace.SpanProcessor {
	var old *sdktrace.SpanProcessor
	if p == nil {
		old = sp.ptr.Swap(nil)
	} else {
		old = sp.ptr.Swap(&p)
	}
	if old == nil {
		return nil
	}

	return *old
}

func (sp *spanProcessor) load() sdktrace.SpanProcessor {
	if p := sp.ptr.Load(); p != nil {
		return *p
	}

	return nil
}

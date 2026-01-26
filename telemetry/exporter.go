package telemetry

import (
	"context"
	"sync/atomic"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type spanExporter struct {
	ptr atomic.Pointer[sdktrace.SpanExporter]
}

func (se *spanExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	if exp := se.load(); exp != nil {
		return exp.ExportSpans(ctx, spans)
	}

	return nil
}

func (se *spanExporter) Shutdown(ctx context.Context) error {
	if exp := se.load(); exp != nil {
		return exp.Shutdown(ctx)
	}

	return nil
}

func (se *spanExporter) swap(exp sdktrace.SpanExporter) sdktrace.SpanExporter {
	var old *sdktrace.SpanExporter
	if exp == nil {
		old = se.ptr.Swap(nil)
	} else {
		old = se.ptr.Swap(&exp)
	}
	if old != nil {
		return *old
	}

	return nil
}

func (se *spanExporter) load() sdktrace.SpanExporter {
	if exp := se.ptr.Load(); exp != nil {
		return *exp
	}

	return nil
}

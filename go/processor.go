package wasi_otel

import (
	"context"
	"fmt"

	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Source: https://github.com/open-telemetry/opentelemetry-go/blob/20cd7f871fce8fae83ea35958775043d6e94b925/sdk/trace/span_processor.go#L15

type WasiProcessor struct {
	IsShutdown bool
}

func NewWasiProcessor() *WasiProcessor {
	return &WasiProcessor{IsShutdown: false}
}

func DefaultWasiProcessor() *WasiProcessor {
	return NewWasiProcessor()
}

func (p WasiProcessor) OnStart(parent context.Context, s sdkTrace.ReadWriteSpan) {
	if p.IsShutdown {
		return
	}

	spanCtx := s.SpanContext()

	// TODO: Not sure if this is correct.
	// Do we want to include the parent context?
	// If not, we need to find another method that doesn't call for the parent context and returns the correct type
	p.OnStart(trace.ContextWithSpanContext(parent, spanCtx), s)
}

func (p WasiProcessor) OnEnd(s sdkTrace.ReadOnlySpan) {
	if p.IsShutdown {
		return
	}

	p.OnEnd(s)
}

func (p WasiProcessor) ForceFlush(ctx context.Context) error {
	if p.IsShutdown {
		return fmt.Errorf("processor already shutdown")
	}

	return nil
}

func (p WasiProcessor) Shutdown(ctx context.Context) error {
	return p.ForceFlush(ctx)
}

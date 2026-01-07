package tracing

import (
	"context"

	wasiTrace "github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_tracing"
	"go.opentelemetry.io/otel/sdk/trace"
	traceApi "go.opentelemetry.io/otel/trace"
)

type WasiSpanProcessor struct {
	trace.SpanProcessor
}

func NewWasiSpanProcessor() WasiSpanProcessor {
	return WasiSpanProcessor{
		SpanProcessor: trace.NewSimpleSpanProcessor(newWasiSpanExporter()),
	}
}

func (w WasiSpanProcessor) OnStart(_parent context.Context, span trace.ReadWriteSpan) {
	wasiTrace.OnStart(toWasiSpanContext(span.SpanContext()))
}

func (w WasiSpanProcessor) OnEnd(s trace.ReadOnlySpan) {
	wasiTrace.OnEnd(toWasiSpanData(s))
}

type wasiSpanExporter struct {
	trace.SpanExporter
}

// A placeholder struct to satisfy the interface
func newWasiSpanExporter() *wasiSpanExporter {
	return &wasiSpanExporter{}
}

type TraceContextPropagator struct{}

func NewTraceContextPropagator() TraceContextPropagator {
	return TraceContextPropagator{}
}

// Retrieves trace context from a WASI host and combines it with the current trace context.
func (t *TraceContextPropagator) Extract(ctx context.Context) context.Context {
	return traceApi.ContextWithRemoteSpanContext(ctx, toOtelSpanContext(wasiTrace.OuterSpanContext()))
}

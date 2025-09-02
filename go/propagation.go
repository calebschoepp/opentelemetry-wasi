package wasi_otel

import (
	"context"

	wasiTrace "github.com/calebschoepp/opentelemetry-wasi/internal/wasi/otel/tracing"
	"go.opentelemetry.io/otel/trace"
)

type WasiPropagator interface {
	Extract(cx context.Context) context.Context
}

type TraceContextPropagator struct{}

func NewTraceContextPropagator() *TraceContextPropagator {
	return &TraceContextPropagator{}
}

func DefaultTraceContextPropagator() *TraceContextPropagator {
	return NewTraceContextPropagator()
}

func (t TraceContextPropagator) Extract(cx context.Context) context.Context {
	// Retrieving span context from the wasm host
	hostCx := wasiTrace.OuterSpanContext()

	// Converting wasm host TraceState type to otel TraceState type
	var traceState trace.TraceState
	for _, entry := range hostCx.TraceState.Slice() {
		ts, err := traceState.Insert(entry[0], entry[1])
		if err != nil {
			// TODO: not sure how to handle this error
		}
		traceState = ts
	}

	cfg := trace.SpanContextConfig{
		TraceID:    trace.TraceID([]byte(hostCx.TraceID)),
		SpanID:     trace.SpanID([]byte(hostCx.SpanID)),
		TraceFlags: trace.TraceFlags(hostCx.TraceFlags),
		TraceState: traceState,
		Remote:     hostCx.IsRemote,
	}

	convertedCx := trace.NewSpanContext(cfg)

	return trace.ContextWithRemoteSpanContext(cx, convertedCx)
}

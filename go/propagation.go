package wasi_otel

import (
	"context"

	wasmTrace "github.com/calebschoepp/opentelemetry-wasi/internal/wasi/otel/tracing"
	otelTrace "go.opentelemetry.io/otel/trace"
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
	hostCx := wasmTrace.OuterSpanContext()

	// Converting wasm host TraceState type to otel TraceState type
	otelTraceState := otelTrace.TraceState{}
	for _, unit := range hostCx.TraceState.Slice() {
		traceState, err := otelTraceState.Insert(unit[0], unit[1])
		if err != nil {
			// TODO: not sure how to handle this error
		}
		otelTraceState = traceState
	}

	cfg := otelTrace.SpanContextConfig{
		TraceID:    otelTrace.TraceID([]byte(hostCx.TraceID)),
		SpanID:     otelTrace.SpanID([]byte(hostCx.SpanID)),
		TraceFlags: otelTrace.TraceFlags(hostCx.TraceFlags),
		TraceState: otelTraceState,
		Remote:     hostCx.IsRemote,
	}

	convertedCx := otelTrace.NewSpanContext(cfg)

	return otelTrace.ContextWithRemoteSpanContext(cx, convertedCx)
}

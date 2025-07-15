package trace

// #cgo CFLAGS: -Wno-unused-parameter -Wno-switch-bool
// #include<tracing.h>
// #include<stdlib.h>
// #include<stdint.h>
import "C"
import (
	"context"
	"fmt"

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
	fmt.Println("Retrieving outer span context...")
	// Retrieving span context from the wasm host
	var cSpanContext C.tracing_span_context_t
	C.outer_span_context(&cSpanContext)
	defer C.tracing_span_context_free(&cSpanContext)
	fmt.Println("Outer span context has been retrieved")

	fmt.Println("Converting trace state...")
	// Converting wasm host TraceState type to otel TraceState type
	otelTraceState := trace.TraceState{}
	for key, value := range otelTraceStateToGoMap(cSpanContext.trace_state) {
		traceState, err := otelTraceState.Insert(key, value)
		if err != nil {
			// TODO: not sure how to handle this error
		}
		otelTraceState = traceState
	}
	fmt.Println("Trace state has been converted")

	fmt.Println("Building span context config...")
	cfg := trace.SpanContextConfig{
		TraceID:    trace.TraceID(otelStringToGoByteSlice(cSpanContext.trace_id)),
		SpanID:     trace.SpanID(otelStringToGoByteSlice(cSpanContext.span_id)),
		TraceFlags: trace.TraceFlags(cSpanContext.trace_flags),
		TraceState: otelTraceState,
		Remote:     bool(cSpanContext.is_remote),
	}
	fmt.Println("Span context config has been built")

	fmt.Println("Converting context to otel context...")
	convertedCx := trace.NewSpanContext(cfg)
	fmt.Println("Context has been converted to otel context")

	return trace.ContextWithRemoteSpanContext(cx, convertedCx)
}

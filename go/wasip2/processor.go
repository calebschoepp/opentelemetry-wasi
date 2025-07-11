package wasi_otel

import (
	"context"
	"fmt"

	wasmTrace "github.com/calebschoepp/opentelemetry-wasi/internal/wasi/otel/tracing"
	"go.bytecodealliance.org/cm"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
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

	spanCx := s.SpanContext()

	traceStateList := cm.NewList[[2]string](nil, 0)
	for k, v := range spanCx.TraceState().Walk {
		traceStateList = cm.NewList(&[2]string{k, v}, 2)
	}

	convertedCx := wasmTrace.SpanContext{
		TraceID:    wasmTrace.TraceID(spanCx.TraceID().String()),
		SpanID:     wasmTrace.SpanID(spanCx.SpanID().String()),
		TraceFlags: wasmTrace.TraceFlags(spanCx.TraceFlags()),
		IsRemote:   spanCx.IsRemote(),
		TraceState: wasmTrace.TraceState(traceStateList),
	}

	wasmTrace.OnStart(convertedCx)
}

func (p WasiProcessor) OnEnd(s sdkTrace.ReadOnlySpan) {
	if p.IsShutdown {
		return
	}

	wasmTrace.OnEnd(wasmTrace.SpanData{
		SpanContext:  wasmTrace.SpanContext{}, // TODO: Fill in
		ParentSpanID: s.Parent().SpanID().String(),
		SpanKind:     wasmTrace.SpanKind(s.SpanKind()),
		Name:         s.Name(),
		StartTime: wasmTrace.DateTime{
			Seconds:     uint64(s.StartTime().Second()),     // TODO: make sure these are returning the actual time in SECONDS, not just the seconds place
			Nanoseconds: uint32(s.StartTime().Nanosecond()), // TODO: same here, but for NANOSECONDS
		},
		EndTime: wasmTrace.DateTime{
			Seconds:     uint64(s.EndTime().Second()),     // TODO: same as above
			Nanoseconds: uint32(s.EndTime().Nanosecond()), // TODO: same as above
		},
		Attributes:           cm.List[wasmTrace.KeyValue]{},    // TODO: Fill in
		Events:               cm.List[wasmTrace.Event]{},       // TODO: Fill in
		Links:                cm.List[wasmTrace.Link]{},        // TODO: Fill in
		Status:               wasmTrace.Status{},               // TODO: Fill in
		InstrumentationScope: wasmTrace.InstrumentationScope{}, // TODO: Fill in
		DroppedAttributes:    0,                                // TODO: Fill in
		DroppedEvents:        0,                                // TODO: Fill in
		DroppedLinks:         0,                                // TODO: Fill in
	})
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

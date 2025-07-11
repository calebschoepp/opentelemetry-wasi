//go:build cgo
// +build cgo

package trace

/*
#include <stdint.h>
#include <stdlib.h>
#include "traces.h"
*/
import "C"
import (
	"context"
	"fmt"

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

	convertedCx := C.traces_span_context_t{
		trace_id:    goStringToOtelString(spanCx.TraceID().String()),
		span_id:     goStringToOtelString(spanCx.SpanID().String()),
		trace_flags: C.uint8_t(spanCx.TraceFlags()),
		is_remote:   C.bool(spanCx.IsRemote()),
		trace_state: *goTraceStateToOtelTraceState(spanCx.TraceState()),
	}

	C.traces_on_start(&convertedCx)
}

func (p WasiProcessor) OnEnd(s sdkTrace.ReadOnlySpan) {
	if p.IsShutdown {
		return
	}

	C.traces_on_end(&C.traces_span_data_t{
		span_context: C.traces_span_context_t{
			trace_id:    goStringToOtelString(s.SpanContext().TraceID().String()),
			span_id:     goStringToOtelString(s.SpanContext().SpanID().String()),
			trace_flags: C.traces_trace_flags_t(s.SpanContext().TraceFlags()),
			is_remote:   C.bool(s.SpanContext().IsRemote()),
			trace_state: *goTraceStateToOtelTraceState(s.SpanContext().TraceState()),
		},
		parent_span_id:        goStringToOtelString(s.Parent().SpanID().String()),
		span_kind:             C.uint8_t(s.SpanKind()),
		name:                  goStringToOtelString(s.Name()),
		start_time:            goTimeToOtelTime(s.StartTime()),
		end_time:              goTimeToOtelTime(s.EndTime()),
		attributes:            *goAttributesToOtelAttributes(s.Attributes()),
		events:                *goEventsToOtelEvents(s.Events()),
		links:                 *goListLinkToOtelListLink(s.Links()),
		status:                goStatusToOtelStatus(s.Status()),
		instrumentation_scope: goInstrumentationScopeToOtelInstrumentationScope(s.InstrumentationScope()),
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

package trace

// #cgo CFLAGS: -Wno-unused-parameter -Wno-switch-bool
// #include<tracing.h>
// #include<stdlib.h>
// #include<stdint.h>
import "C"
import (
	"context"
	"fmt"
	"sync/atomic"

	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

// Source: https://github.com/open-telemetry/opentelemetry-go/blob/20cd7f871fce8fae83ea35958775043d6e94b925/sdk/trace/span_processor.go#L15

type WasiProcessor struct {
	IsShutdown atomic.Bool
}

func NewWasiProcessor() *WasiProcessor {
	return &WasiProcessor{}
}

func DefaultWasiProcessor() *WasiProcessor {
	return NewWasiProcessor()
}

func (p *WasiProcessor) OnStart(parent context.Context, s sdkTrace.ReadWriteSpan) {
	if p.IsShutdown.Load() {
		return
	}

	spanCx := s.SpanContext()

	fmt.Println("Converting span context to `traces_span_context...`")
	convertedCx := &C.tracing_span_context_t{
		trace_id:    goStringToOtelString(spanCx.TraceID().String()),
		span_id:     goStringToOtelString(spanCx.SpanID().String()),
		trace_flags: C.uint8_t(spanCx.TraceFlags()),
		is_remote:   C._Bool(spanCx.IsRemote()),
		trace_state: *goTraceStateToOtelTraceState(spanCx.TraceState()),
	}
	defer C.tracing_span_context_free(convertedCx)
	fmt.Println("Successfully converted to `traces_span_context`")

	C.on_start(convertedCx)
}

func (p *WasiProcessor) OnEnd(s sdkTrace.ReadOnlySpan) {
	if p.IsShutdown.Load() {
		return
	}

	fmt.Println("Calling `traces_on_end`...")
	data := &C.tracing_span_data_t{
		span_context: C.tracing_span_context_t{
			trace_id:    goStringToOtelString(s.SpanContext().TraceID().String()),
			span_id:     goStringToOtelString(s.SpanContext().SpanID().String()),
			trace_flags: C.tracing_trace_flags_t(s.SpanContext().TraceFlags()),
			is_remote:   C._Bool(s.SpanContext().IsRemote()),
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
	}
	defer C.tracing_span_data_free(data)

	C.on_end(data)

	fmt.Println("Success calling `traces_on_end`")
}

func (p *WasiProcessor) ForceFlush(ctx context.Context) error {
	if p.IsShutdown.Load() {
		return fmt.Errorf("processor already shutdown")
	}

	return nil
}

func (p *WasiProcessor) Shutdown(ctx context.Context) error {
	err := p.ForceFlush(ctx)
	if p.IsShutdown.Swap(true) {
		return err
	}

	return nil
}

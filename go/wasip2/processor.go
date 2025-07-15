package wasi_otel

import (
	"context"
	"fmt"
	"sync/atomic"

	wasiTrace "github.com/calebschoepp/opentelemetry-wasi/internal/wasi/otel/tracing"
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
	// TODO: `parent` seems never to be used...
	if p.IsShutdown.Load() {
		return
	}

	wasiTrace.OnStart(toWasiSpanContext(s.SpanContext()))
}

func (p *WasiProcessor) OnEnd(s sdkTrace.ReadOnlySpan) {
	if p.IsShutdown.Load() {
		return
	}

	wasiTrace.OnEnd(wasiTrace.SpanData{
		SpanContext:          toWasiSpanContext(s.SpanContext()),
		ParentSpanID:         s.Parent().SpanID().String(),
		SpanKind:             wasiTrace.SpanKind(s.SpanKind()),
		Name:                 s.Name(),
		StartTime:            toWasiDateTime(s.StartTime()),
		EndTime:              toWasiDateTime(s.EndTime()),
		Attributes:           toWasiAttributes(s.Attributes()),
		Events:               toWasiEvents(s.Events()),
		Links:                toWasiLinks(s.Links()),
		Status:               toWasiStatus(s.Status()),
		InstrumentationScope: toWasiInstrumentationScope(s.InstrumentationScope()),
		DroppedAttributes:    uint32(s.DroppedAttributes()),
		DroppedEvents:        uint32(s.DroppedEvents()),
		DroppedLinks:         uint32(s.DroppedLinks()),
	})
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

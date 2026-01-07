package tracing

import (
	"encoding/hex"
	"fmt"

	"github.com/calebschoepp/opentelemetry-wasi/types"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_tracing"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wit_types"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/trace"
	traceApi "go.opentelemetry.io/otel/trace"
)

func toWasiSpanData(s trace.ReadOnlySpan) wasi_otel_tracing.SpanData {
	return wasi_otel_tracing.SpanData{
		SpanContext:          toWasiSpanContext(s.SpanContext()),
		ParentSpanId:         s.Parent().SpanID().String(),
		SpanKind:             toWasiSpanKind(s.SpanKind()),
		Name:                 s.Name(),
		StartTime:            types.ToWasiTime(s.StartTime()),
		EndTime:              types.ToWasiTime(s.EndTime()),
		Attributes:           types.ToWasiAttributes(s.Attributes()),
		Events:               toWasiEvents(s.Events()),
		Links:                toWasiLinks(s.Links()),
		Status:               toWasiStatus(s.Status()),
		InstrumentationScope: types.ToWasiInstrumentationScope(s.InstrumentationScope()),
		DroppedAttributes:    uint32(s.DroppedAttributes()),
		DroppedEvents:        uint32(s.DroppedEvents()),
		DroppedLinks:         uint32(s.DroppedLinks()),
	}
}

func toWasiStatus(s trace.Status) wasi_otel_tracing.Status {
	switch s.Code {
	case codes.Unset:
		return wasi_otel_tracing.MakeStatusUnset()
	case codes.Error:
		return wasi_otel_tracing.MakeStatusError(s.Description)
	case codes.Ok:
		return wasi_otel_tracing.MakeStatusOk()
	default:
		return wasi_otel_tracing.MakeStatusUnset()
	}
}

func toWasiLinks(links []trace.Link) []wasi_otel_tracing.Link {
	result := make([]wasi_otel_tracing.Link, len(links))
	for i, l := range links {
		result[i] = wasi_otel_tracing.Link{
			SpanContext: toWasiSpanContext(l.SpanContext),
			Attributes:  types.ToWasiAttributes(l.Attributes),
		}
	}

	return result
}

func toWasiEvents(events []trace.Event) []wasi_otel_tracing.Event {
	result := make([]wasi_otel_tracing.Event, len(events))
	for i, e := range events {
		result[i] = wasi_otel_tracing.Event{
			Name:       e.Name,
			Time:       types.ToWasiTime(e.Time),
			Attributes: types.ToWasiAttributes(e.Attributes),
		}
	}

	return result
}

func toWasiSpanKind(sk traceApi.SpanKind) wasi_otel_tracing.SpanKind {
	switch sk {
	case traceApi.SpanKindClient:
		return wasi_otel_tracing.SpanKindClient
	case traceApi.SpanKindConsumer:
		return wasi_otel_tracing.SpanKindConsumer
	case traceApi.SpanKindInternal:
		return wasi_otel_tracing.SpanKindInternal
	case traceApi.SpanKindProducer:
		return wasi_otel_tracing.SpanKindProducer
	case traceApi.SpanKindServer:
		return wasi_otel_tracing.SpanKindServer
	case traceApi.SpanKindUnspecified:
		panic("SpanKindUnspecified is not implemented")
	default:
		panic("unimplemented type")
	}
}

func toOtelSpanContext(ctx wasi_otel_tracing.SpanContext) traceApi.SpanContext {
	tid, err := hex.DecodeString(ctx.TraceId)
	if err != nil {
		panic(fmt.Sprintf("invalid trace ID: %v", err))
	}
	if len(tid) != 16 {
		panic(fmt.Sprintf("trace ID must be 16 bytes, got %d", len(tid)))
	}

	sid, err := hex.DecodeString(ctx.SpanId)
	if err != nil {
		panic(fmt.Sprintf("invalid span ID: %v", err))
	}
	if len(sid) != 8 {
		panic(fmt.Sprintf("span id must be 8 bytes, got %d", len(sid)))
	}

	traceState := traceApi.TraceState{}
	for _, kv := range ctx.TraceState {
		ts, err := traceState.Insert(kv.F0, kv.F1)
		if err != nil {
			panic(fmt.Sprintf("invalid trace state entry %s=%s: %v", kv.F0, kv.F1, err))
		}
		traceState = ts
	}

	return traceApi.NewSpanContext(traceApi.SpanContextConfig{
		TraceID:    [16]byte(tid),
		SpanID:     [8]byte(sid),
		TraceFlags: traceApi.FlagsSampled,
		TraceState: traceState,
		Remote:     ctx.IsRemote,
	})
}

func toWasiSpanContext(s traceApi.SpanContext) wasi_otel_tracing.SpanContext {
	traceState := make([]wit_types.Tuple2[string, string], s.TraceState().Len())
	s.TraceState().Walk(func(key, value string) bool {
		traceState = append(traceState, wit_types.Tuple2[string, string]{
			F0: key,
			F1: value,
		})
		return true
	})

	return wasi_otel_tracing.SpanContext{
		TraceId:    s.TraceID().String(),
		SpanId:     s.SpanID().String(),
		TraceFlags: wasi_otel_tracing.TraceFlagsSampled,
		IsRemote:   s.IsRemote(),
		TraceState: traceState,
	}
}

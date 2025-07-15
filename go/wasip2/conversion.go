package wasi_otel

import (
	"time"

	wasiTrace "github.com/calebschoepp/opentelemetry-wasi/internal/wasi/otel/tracing"
	"go.bytecodealliance.org/cm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func toWasiTraceState(s trace.TraceState) wasiTrace.TraceState {
	// TODO: I don't trust that this is the proper way to create/add to a cm.List;
	// it seems like it would just overwrite the existing list with a single, latest entry
	result := cm.NewList[[2]string](nil, 0)
	s.Walk(func(k, v string) bool {
		result = cm.NewList(&[2]string{k, v}, 2)
		return true
	})

	return wasiTrace.TraceState(result)
}

func toWasiDateTime(t time.Time) wasiTrace.DateTime {
	// TODO: Make sure this is retrieving the correct time
	return wasiTrace.DateTime{
		Seconds:     uint64(t.Second()),
		Nanoseconds: uint32(t.Nanosecond()),
	}
}

func toWasiAttributes(attributes []attribute.KeyValue) cm.List[wasiTrace.KeyValue] {
	result := cm.NewList[wasiTrace.KeyValue](nil, 0)
	for _, entry := range attributes {
		result = cm.NewList(&wasiTrace.KeyValue{
			Key:   wasiTrace.Key(entry.Key),
			Value: toWasiValue(entry.Value),
		}, len(attributes))
	}

	return result
}

func toWasiEvents(events []sdkTrace.Event) cm.List[wasiTrace.Event] {
	result := cm.NewList[wasiTrace.Event](nil, 0)
	for _, event := range events {
		result = cm.NewList(&wasiTrace.Event{
			Name:       event.Name,
			Time:       toWasiDateTime(event.Time),
			Attributes: toWasiAttributes(event.Attributes),
		}, len(events))
	}

	return result
}

func toWasiSpanContext(s trace.SpanContext) wasiTrace.SpanContext {
	return wasiTrace.SpanContext{
		TraceID:    wasiTrace.TraceID(s.TraceID().String()),
		SpanID:     wasiTrace.SpanID(s.SpanID().String()),
		TraceFlags: wasiTrace.TraceFlags(s.TraceFlags()),
		IsRemote:   s.IsRemote(),
		TraceState: toWasiTraceState(s.TraceState()),
	}
}

func toWasiLinks(links []sdkTrace.Link) cm.List[wasiTrace.Link] {
	result := cm.NewList[wasiTrace.Link](nil, 0)
	for _, link := range links {
		result = cm.NewList(&wasiTrace.Link{
			SpanContext: toWasiSpanContext(link.SpanContext),
			Attributes:  toWasiAttributes(link.Attributes),
		}, len(links))
	}

	return result
}

func toWasiInstrumentationScope(s instrumentation.Scope) wasiTrace.InstrumentationScope {
	return wasiTrace.InstrumentationScope{
		Name:       s.Name,
		Version:    getOption(s.Version),
		SchemaURL:  getOption(s.SchemaURL),
		Attributes: toWasiAttributes(s.Attributes.ToSlice()),
	}
}

func toWasiStatus(s sdkTrace.Status) wasiTrace.Status {
	switch s.Code {
	case codes.Error:
		return wasiTrace.StatusError(s.Description)
	case codes.Ok:
		return wasiTrace.StatusOK()
	default:
		return wasiTrace.StatusUnset()
	}
}

func toWasiValue(v attribute.Value) wasiTrace.Value {
	switch v.Type() {
	case attribute.STRING:
		return wasiTrace.ValueString_(v.AsString())
	case attribute.BOOL:
		return wasiTrace.ValueBool(v.AsBool())
	case attribute.FLOAT64:
		return wasiTrace.ValueF64(v.AsFloat64())
	case attribute.INT64:
		return wasiTrace.ValueS64(v.AsInt64())
	case attribute.STRINGSLICE:
		stringSlice := v.AsStringSlice()
		return wasiTrace.ValueStringArray(cm.ToList(stringSlice))
	case attribute.BOOLSLICE:
		boolSlice := v.AsBoolSlice()
		return wasiTrace.ValueBoolArray(cm.ToList(boolSlice))
	case attribute.FLOAT64SLICE:
		floatSlice := v.AsFloat64Slice()
		return wasiTrace.ValueF64Array(cm.ToList(floatSlice))
	case attribute.INT64SLICE:
		intSlice := v.AsInt64Slice()
		return wasiTrace.ValueS64Array(cm.ToList(intSlice))
	default:
		return wasiTrace.ValueString_(v.AsString())
	}
}

func getOption[T comparable](v T) cm.Option[T] {
	// TODO: Claude claims that this is the best way to do this function. VALIDATE
	var empty T
	if v == empty {
		return cm.None[T]()
	}
	return cm.Some(v)
}

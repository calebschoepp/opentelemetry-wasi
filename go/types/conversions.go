package types

import (
	"encoding/json"
	"time"

	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_clocks_wall_clock"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_types"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wit_types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/resource"
)

func ToWasiResource(r resource.Resource) wasi_otel_types.Resource {
	return wasi_otel_types.Resource{
		Attributes: ToWasiAttributes(r.Attributes()),
		SchemaUrl:  ToWasiOptStr(r.SchemaURL()),
	}
}

func ToWasiAttributes(attrs []attribute.KeyValue) []wasi_otel_types.KeyValue {
	result := make([]wasi_otel_types.KeyValue, len(attrs))
	for i, attr := range attrs {
		result[i] = wasi_otel_types.KeyValue{
			Key:   string(attr.Key),
			Value: otelValueToJson(attr.Value),
		}

	}

	return result
}

func otelValueToJson(v attribute.Value) string {
	switch v.Type() {
	case attribute.BOOL:
		bytes, _ := json.Marshal(v.AsBool())
		return string(bytes)
	case attribute.BOOLSLICE:
		bytes, _ := json.Marshal(v.AsBoolSlice())
		return string(bytes)
	case attribute.INT64:
		bytes, _ := json.Marshal(v.AsInt64())
		return string(bytes)
	case attribute.INT64SLICE:
		bytes, _ := json.Marshal(v.AsInt64Slice())
		return string(bytes)
	case attribute.FLOAT64:
		bytes, _ := json.Marshal(v.AsFloat64())
		return string(bytes)
	case attribute.FLOAT64SLICE:
		bytes, _ := json.Marshal(v.AsFloat64Slice())
		return string(bytes)
	case attribute.STRING:
		bytes, _ := json.Marshal(v.AsString())
		return string(bytes)
	case attribute.STRINGSLICE:
		bytes, _ := json.Marshal(v.AsStringSlice())
		return string(bytes)
	case attribute.INVALID:
		panic("invalid type")
	default:
		panic("unsupported type")
	}
}

func ToWasiOptStr(s string) wit_types.Option[string] {
	if s == "" {
		return wit_types.None[string]()
	}

	return wit_types.Some(s)
}

func ToWasiInstrumentationScope(s instrumentation.Scope) wasi_otel_types.InstrumentationScope {
	return wasi_otel_types.InstrumentationScope{
		Name:       s.Name,
		Version:    ToWasiOptStr(s.Version),
		SchemaUrl:  ToWasiOptStr(s.SchemaURL),
		Attributes: ToWasiAttributes(s.Attributes.ToSlice()),
	}
}

func ToWasiTime(t time.Time) wasi_clocks_wall_clock.Datetime {
	return wasi_clocks_wall_clock.Datetime{
		Seconds:     uint64(t.Unix()),
		Nanoseconds: uint32(t.Nanosecond()),
	}
}

package logs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/calebschoepp/opentelemetry-wasi/types"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_clocks_wall_clock"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_logs"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_tracing"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_types"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wit_types"
	logApi "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
)

func toWasiLogRecord(r log.Record) wasi_otel_logs.LogRecord {
	var ts wit_types.Option[wasi_otel_logs.Datetime]
	if r.Timestamp().IsZero() {
		ts = wit_types.None[wasi_clocks_wall_clock.Datetime]()
	} else {
		ts = wit_types.Some(types.ToWasiTime(r.Timestamp()))
	}

	var ots wit_types.Option[wasi_otel_logs.Datetime]
	if r.ObservedTimestamp().IsZero() {
		ots = wit_types.None[wasi_clocks_wall_clock.Datetime]()
	} else {
		ots = wit_types.Some(types.ToWasiTime(r.ObservedTimestamp()))
	}

	var sn wit_types.Option[uint8]
	if r.Severity() == logApi.SeverityUndefined {
		sn = wit_types.None[uint8]()
	} else {
		sn = wit_types.Some(uint8(r.Severity()))
	}

	var st wit_types.Option[string]
	if r.SeverityText() == "" {
		st = wit_types.None[string]()
	} else {
		st = wit_types.Some(r.SeverityText())
	}

	var attrs wit_types.Option[[]wasi_otel_types.KeyValue]
	if r.AttributesLen() == 0 {
		attrs = wit_types.None[[]wasi_otel_types.KeyValue]()
	} else {
		attrList := make([]wasi_otel_types.KeyValue, 0)
		r.WalkAttributes(func(attr logApi.KeyValue) bool {
			attrList = append(attrList, wasi_otel_types.KeyValue{
				Key:   attr.Key,
				Value: OtelLogValueToJson(attr.Value),
			})

			return true
		})

		attrs = wit_types.Some(attrList)
	}

	var res wit_types.Option[wasi_otel_types.Resource]
	if r.Resource() == nil {
		res = wit_types.None[wasi_otel_logs.Resource]()
	} else {
		res = wit_types.Some(types.ToWasiResource(*r.Resource()))
	}

	var is wit_types.Option[wasi_otel_types.InstrumentationScope]
	if r.InstrumentationScope().Name == "" {
		is = wit_types.None[wasi_otel_logs.InstrumentationScope]()
	} else {
		is = wit_types.Some(types.ToWasiInstrumentationScope(r.InstrumentationScope()))
	}

	var tf wit_types.Option[wasi_otel_tracing.TraceFlags]
	if !r.TraceFlags().IsSampled() {
		tf = wit_types.None[wasi_otel_tracing.TraceFlags]()
	} else {
		tf = wit_types.Some(wasi_otel_tracing.TraceFlagsSampled)
	}

	return wasi_otel_logs.LogRecord{
		Timestamp:            ts,
		ObservedTimestamp:    ots,
		SeverityNumber:       sn,
		SeverityText:         st,
		Body:                 types.ToWasiOptStr(OtelLogValueToJson(r.Body())),
		Attributes:           attrs,
		EventName:            types.ToWasiOptStr(r.EventName()),
		Resource:             res,
		InstrumentationScope: is,
		TraceId:              types.ToWasiOptStr(r.TraceID().String()),
		SpanId:               types.ToWasiOptStr(r.SpanID().String()),
		TraceFlags:           tf,
	}
}

func OtelLogValueToJson(v logApi.Value) string {
	switch v.Kind() {
	case logApi.KindBool:
		bytes, err := json.Marshal(v.AsBool())
		if err != nil {
			panic(err)
		}
		return string(bytes)
	case logApi.KindFloat64:
		bytes, err := json.Marshal(v.AsFloat64())
		if err != nil {
			panic(err)
		}
		return string(bytes)
	case logApi.KindInt64:
		bytes, err := json.Marshal(v.AsInt64())
		if err != nil {
			panic(err)
		}
		return string(bytes)
	case logApi.KindEmpty:
		bytes, err := json.Marshal("")
		if err != nil {
			panic(err)
		}
		return string(bytes)
	case logApi.KindBytes:
		bytes, err := json.Marshal(fmt.Sprintf("{base64}:%s", base64.StdEncoding.EncodeToString(v.AsBytes())))
		if err != nil {
			panic(err)
		}
		return string(bytes)
	case logApi.KindString:
		bytes, err := json.Marshal(v.AsString())
		if err != nil {
			panic(err)
		}
		return string(bytes)
	case logApi.KindSlice:
		slice := v.AsSlice()
		result := make([]any, len(slice))
		for i, item := range slice {
			var temp any
			if err := json.Unmarshal([]byte(OtelLogValueToJson(item)), &temp); err != nil {
				panic(err)
			}
			result[i] = temp
		}
		bytes, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		return string(bytes)
	case logApi.KindMap:
		kvs := v.AsMap()
		result := make(map[string]any)
		for _, kv := range kvs {
			var temp any
			if err := json.Unmarshal([]byte(OtelLogValueToJson(kv.Value)), &temp); err != nil {
				panic(err)
			}
			result[kv.Key] = temp
		}
		bytes, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		return string(bytes)
	default:
		panic("unsupported type")
	}
}

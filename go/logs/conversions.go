package logs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/calebschoepp/opentelemetry-wasi/internal/wasi_clocks_wall_clock"
	"github.com/calebschoepp/opentelemetry-wasi/internal/wasi_otel_logs"
	"github.com/calebschoepp/opentelemetry-wasi/internal/wasi_otel_tracing"
	"github.com/calebschoepp/opentelemetry-wasi/internal/wasi_otel_types"
	"github.com/calebschoepp/opentelemetry-wasi/types"
	witTypes "go.bytecodealliance.org/pkg/wit/types"
	logApi "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/log"
)

func toWasiLogRecord(r log.Record) wasi_otel_logs.LogRecord {
	var ts witTypes.Option[wasi_otel_logs.Datetime]
	if r.Timestamp().IsZero() {
		ts = witTypes.None[wasi_clocks_wall_clock.Datetime]()
	} else {
		ts = witTypes.Some(types.ToWasiTime(r.Timestamp()))
	}

	var ots witTypes.Option[wasi_otel_logs.Datetime]
	if r.ObservedTimestamp().IsZero() {
		ots = witTypes.None[wasi_clocks_wall_clock.Datetime]()
	} else {
		ots = witTypes.Some(types.ToWasiTime(r.ObservedTimestamp()))
	}

	var sn witTypes.Option[uint8]
	if r.Severity() == logApi.SeverityUndefined {
		sn = witTypes.None[uint8]()
	} else {
		sn = witTypes.Some(uint8(r.Severity()))
	}

	var st witTypes.Option[string]
	if r.SeverityText() == "" {
		st = witTypes.None[string]()
	} else {
		st = witTypes.Some(r.SeverityText())
	}

	var attrs witTypes.Option[[]wasi_otel_types.KeyValue]
	if r.AttributesLen() == 0 {
		attrs = witTypes.None[[]wasi_otel_types.KeyValue]()
	} else {
		attrList := make([]wasi_otel_types.KeyValue, 0)
		r.WalkAttributes(func(attr logApi.KeyValue) bool {
			attrList = append(attrList, wasi_otel_types.KeyValue{
				Key:   attr.Key,
				Value: OtelLogValueToJson(attr.Value),
			})

			return true
		})

		attrs = witTypes.Some(attrList)
	}

	var res witTypes.Option[wasi_otel_types.Resource]
	if r.Resource() == nil {
		res = witTypes.None[wasi_otel_logs.Resource]()
	} else {
		res = witTypes.Some(types.ToWasiResource(*r.Resource()))
	}

	var is witTypes.Option[wasi_otel_types.InstrumentationScope]
	if r.InstrumentationScope().Name == "" {
		is = witTypes.None[wasi_otel_logs.InstrumentationScope]()
	} else {
		is = witTypes.Some(types.ToWasiInstrumentationScope(r.InstrumentationScope()))
	}

	var tf witTypes.Option[wasi_otel_tracing.TraceFlags]
	if !r.TraceFlags().IsSampled() {
		tf = witTypes.None[wasi_otel_tracing.TraceFlags]()
	} else {
		tf = witTypes.Some(wasi_otel_tracing.TraceFlagsSampled)
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
		bytes, err := json.Marshal(fmt.Sprintf("data:application/octet-stream;base64,%s", base64.StdEncoding.EncodeToString(v.AsBytes())))
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

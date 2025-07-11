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
	"time"
	"unsafe"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// Translates a `C.traces_trace_state_t` to a map[string]string
func otelTraceStateToGoMap(traceState C.traces_trace_state_t) map[string]string {
	result := make(map[string]string)

	if traceState.ptr == nil || traceState.len == 0 {
		return result
	}

	// Convert C array to Go slice for easier iteration
	length := int(traceState.len)
	kvPairs := (*[1 << 28]C.traces_tuple2_string_string_t)(unsafe.Pointer(traceState.ptr))[:length:length] // TODO: figure out what this does

	for i := 0; i < length; i++ {
		key := otelStringToGoString(kvPairs[i].f0)
		value := otelStringToGoString(kvPairs[i].f1)
		result[key] = value
	}

	return result
}

// Translates from `time.Time` to `C.traces_datetime_t`
func goTimeToOtelTime(t time.Time) C.traces_datetime_t {
	return C.traces_datetime_t{
		seconds:     C.uint64_t(t.Unix()),
		nanoseconds: C.uint32_t(t.Nanosecond()),
	}
}

// Translates from an otel `TraceState` to a `C.traces_trace_state_t`. IMPORTANT: don't forget to free the pointer to the C array
func goTraceStateToOtelTraceState(t trace.TraceState) *C.traces_trace_state_t {
	// Allocate tuple array
	// tupleSize := C.size_t(unsafe.Sizeof(C.traces_tuple2_string_string_t{}))
	// arraySize := C.size_t(t.Len()) * tupleSize
	// cArray := (*C.traces_tuple2_string_string_t)(C.malloc(arraySize))
	cArray := allocateCArray[C.traces_tuple2_string_string_t](t.Len())

	i := 0
	for k, v := range t.Walk {
		// Get pointer to the i-th tuple in the array
		tuplePtr := (*C.traces_tuple2_string_string_t)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(C.traces_tuple2_string_string_t{}),
		))

		// Fill the tuple
		tuplePtr.f0 = goStringToOtelString(k)
		tuplePtr.f1 = goStringToOtelString(v)

		i++
	}

	return &C.traces_trace_state_t{
		ptr: cArray,
		len: C.size_t(t.Len()),
	}
}

// Translates from otel `[]attribute.KeyValue` to `C.traces_list_key_value_t`. IMPORTANT: don't forget to free the pointer to the C array
func goAttributesToOtelAttributes(attributes []attribute.KeyValue) *C.traces_list_key_value_t {
	// entrySize := C.size_t(unsafe.Sizeof(C.traces_key_value_t{}))
	// arraySize := C.size_t(len(attributes)) * entrySize
	// cArray := (*C.traces_key_value_t)(C.malloc(arraySize))
	cArray := allocateCArray[C.traces_key_value_t](len(attributes)) // TODO: Make sure this actually works, and refactor other functions to use this if it does

	for i, attribute := range attributes {
		entryPtr := (*C.traces_key_value_t)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(C.traces_key_value_t{}),
		))

		entryPtr.key = goStringToOtelString(string(attribute.Key))
		entryPtr.value = goValueToOtelValue(attribute.Value)
	}

	return &C.traces_list_key_value_t{
		ptr: cArray,
		len: C.size_t(len(attributes)),
	}
}

func goEventsToOtelEvents(events []sdkTrace.Event) *C.traces_list_event_t {
	cArray := allocateCArray[C.traces_event_t](len(events))
	for i, event := range events {
		entryPtr := (*C.traces_event_t)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(C.traces_event_t{}),
		))

		entryPtr.name = goStringToOtelString(event.Name)
		entryPtr.time = goTimeToOtelTime(event.Time)
		entryPtr.attributes = *goAttributesToOtelAttributes(event.Attributes)
	}

	return &C.traces_list_event_t{
		ptr: cArray,
		len: C.size_t(len(events)),
	}
}

func goListLinkToOtelListLink(links []sdkTrace.Link) *C.traces_list_link_t {
	cArray := allocateCArray[C.traces_link_t](len(links))
	for i, link := range links {
		entryPtr := (*C.traces_link_t)(unsafe.Pointer(
			uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(C.traces_link_t{}),
		))

		entryPtr.span_context = C.traces_span_context_t{
			trace_id:    goStringToOtelString(link.SpanContext.TraceID().String()),
			span_id:     goStringToOtelString(link.SpanContext.SpanID().String()),
			trace_flags: C.traces_trace_flags_t(link.SpanContext.TraceFlags()),
			is_remote:   C.bool(link.SpanContext.IsRemote()),
			trace_state: *goTraceStateToOtelTraceState(link.SpanContext.TraceState()),
		}

		entryPtr.attributes = *goAttributesToOtelAttributes(link.Attributes)
	}

	return &C.traces_list_link_t{
		ptr: cArray,
		len: C.size_t(len(links)),
	}
}

func goStatusToOtelStatus(status sdkTrace.Status) C.traces_status_t {
	var otelStatus C.traces_status_t

	switch status.Code {
	case 1:
		otelStatus.tag = C.TRACES_STATUS_ERROR

		// Error status requires a message
		var errMsg C.traces_string_t
		if len(status.Description) > 0 {
			errMsg = goStringToOtelString(status.Description)
		} else {
			errMsg = goStringToOtelString("unknown error")
		}

		*(*C.traces_string_t)(unsafe.Pointer(&otelStatus.val)) = errMsg

	case 2:
		otelStatus.tag = C.TRACES_STATUS_OK
		// No need to set `val` for OK status

	default:
		otelStatus.tag = C.TRACES_STATUS_UNSET
		// No need to set `val` for UNSET status

	}

	return otelStatus
}

func goValueToOtelValue(value attribute.Value) C.traces_value_t {
	var otelValue C.traces_value_t
	switch value.Type() {
	case attribute.STRING:
		otelValue.tag = C.TRACES_VALUE_STRING
		*(*C.traces_string_t)(unsafe.Pointer(&otelValue.val)) = goStringToOtelString(value.AsString())
	case attribute.BOOL:
		otelValue.tag = C.TRACES_VALUE_BOOL
		*(*C.bool)(unsafe.Pointer(&otelValue.val)) = goBoolToCBool(value.AsBool())
	case attribute.FLOAT64:
		otelValue.tag = C.TRACES_VALUE_F64
		*(*C.double)(unsafe.Pointer(&otelValue.val)) = C.double(value.AsFloat64())
	case attribute.INT64:
		otelValue.tag = C.TRACES_VALUE_S64
		*(*C.int64_t)(unsafe.Pointer(&otelValue.val)) = C.int64_t(value.AsInt64())
	case attribute.STRINGSLICE:
		otelValue.tag = C.TRACES_VALUE_STRING_ARRAY
		*(*C.traces_list_string_t)(unsafe.Pointer(&otelValue.val)) = goStringSliceToOtelListString(value.AsStringSlice())
	case attribute.BOOLSLICE:
		otelValue.tag = C.TRACES_VALUE_BOOL_ARRAY
		*(*C.traces_list_bool_t)(unsafe.Pointer(&otelValue.val)) = goBoolSliceToOtelListBool(value.AsBoolSlice())
	case attribute.FLOAT64SLICE:
		otelValue.tag = C.TRACES_VALUE_F64_ARRAY
		*(*C.traces_list_float64_t)(unsafe.Pointer(&otelValue.val)) = goSliceF64toOtelListF64(value.AsFloat64Slice())
	case attribute.INT64SLICE:
		otelValue.tag = C.TRACES_VALUE_S64_ARRAY
		*(*C.traces_list_s64_t)(unsafe.Pointer(&otelValue.val)) = goSliceS64ToOtelListS64(value.AsInt64Slice())
	default:
		// TODO: unsure how to handle this
	}

	return otelValue
}

func goInstrumentationScopeToOtelInstrumentationScope(scope instrumentation.Scope) C.traces_instrumentation_scope_t {
	return C.traces_instrumentation_scope_t{
		name:       goStringToOtelString(scope.Name),
		version:    goStringToOtelOptionString(scope.Version),
		schema_url: goStringToOtelOptionString(scope.SchemaURL),
		attributes: *goAttributesToOtelAttributes(scope.Attributes.ToSlice()),
	}
}

// TODO: I think the return type may be incorrect
// Allocates a C array of type T and returns a pointer to the array
func allocateCArray[T any](length int) *T {
	var entryType T
	entrySize := C.size_t(unsafe.Sizeof(entryType))
	arraySize := C.size_t(length) * entrySize
	return (*T)(C.malloc(arraySize))
}

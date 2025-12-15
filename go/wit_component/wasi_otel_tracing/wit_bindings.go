package wasi_otel_tracing

import (
        "wit_component/wasi_clocks_wall_clock"
"wit_component/wasi_otel_types"
"wit_component/wit_runtime"
"wit_component/wit_types"
"unsafe"
"runtime"
)

type Datetime = wasi_clocks_wall_clock.Datetime
type KeyValue = wasi_otel_types.KeyValue
type InstrumentationScope = wasi_otel_types.InstrumentationScope
// The trace that this `span-context` belongs to.
// 
// 16 bytes encoded as a hexadecimal string.
type TraceId = string
// The id of this `span-context`.
// 
// 8 bytes encoded as a hexadecimal string.
type SpanId = string

const (
// Whether the `span` should be sampled or not.
TraceFlagsSampled uint8 = 1 << 0
)

// Flags that can be set on a `span-context`.
type TraceFlags = uint8
// Carries system-specific configuration data, represented as a list of key-value pairs. `trace-state` allows multiple tracing systems to participate in the same trace.
// 
// If any invalid keys or values are provided then the `trace-state` will be treated as an empty list.
type TraceState = []wit_types.Tuple2[string, string]

// Identifying trace information about a span that can be serialized and propagated.
type SpanContext struct {
        // The `trace-id` for this `span-context`.
TraceId string
// The `span-id` for this `span-context`.
SpanId string
// The `trace-flags` for this `span-context`.
TraceFlags TraceFlags
// Whether this `span-context` was propagated from a remote parent.
IsRemote bool
// The `trace-state` for this `span-context`.
TraceState []wit_types.Tuple2[string, string] 
}

const (
        // Indicates that the span describes a request to some remote service. This span is usually the parent of a remote server span and does not end until the response is received.
SpanKindClient uint8 = 0
// Indicates that the span covers server-side handling of a synchronous RPC or other remote request. This span is often the child of a remote client span that was expected to wait for a response.
SpanKindServer uint8 = 1
// Indicates that the span describes the initiators of an asynchronous request. This parent span will often end before the corresponding child consumer span, possibly even before the child span starts. In messaging scenarios with batching, tracing individual messages requires a new producer span per message to be created.
SpanKindProducer uint8 = 2
// Indicates that the span describes a child of an asynchronous consumer request.
SpanKindConsumer uint8 = 3
// Default value. Indicates that the span represents an internal operation within an application, as opposed to an operations with remote parents or children.
SpanKindInternal uint8 = 4
)
// Describes the relationship between the Span, its parents, and its children in a trace.
type SpanKind = uint8

// An event describing a specific moment in time on a span and associated attributes.
type Event struct {
        // Event name.
Name string
// Event time.
Time wasi_clocks_wall_clock.Datetime
// Event attributes.
Attributes []wasi_otel_types.KeyValue 
}

// Describes a relationship to another `span`.
type Link struct {
        // Denotes which `span` to link to.
SpanContext SpanContext
// Attributes describing the link.
Attributes []wasi_otel_types.KeyValue 
}

const (
// The default status.
StatusUnset uint8 = 0
// The operation has been validated by an Application developer or Operator to have completed successfully.
StatusOk uint8 = 1
// The operation contains an error with a description.
StatusError uint8 = 2
)

// The `status` of a `span`.
type Status struct {
        tag uint8
        value any
}

func (self Status) Tag() uint8 {
        return self.tag
}

func (self Status) Error() string {
        if self.tag != StatusError {
                panic("tag mismatch")
        }
        return self.value.(string)
}

func MakeStatusUnset() Status {
        return Status{StatusUnset, nil}
}
func MakeStatusOk() Status {
        return Status{StatusOk, nil}
}
func MakeStatusError(value string) Status {
        return Status{StatusError, value}
}



// The data associated with a span.
type SpanData struct {
        // Span context.
SpanContext SpanContext
// Span parent id.
ParentSpanId string
// Span kind.
SpanKind SpanKind
// Span name.
Name string
// Span start time.
StartTime wasi_clocks_wall_clock.Datetime
// Span end time.
EndTime wasi_clocks_wall_clock.Datetime
// Span attributes.
Attributes []wasi_otel_types.KeyValue
// Span events.
Events []Event
// Span Links.
Links []Link
// Span status.
Status Status
// Instrumentation scope that produced this span.
InstrumentationScope wasi_otel_types.InstrumentationScope
// Number of attributes dropped by the span due to limits being reached.
DroppedAttributes uint32
// Number of events dropped by the span due to limits being reached.
DroppedEvents uint32
// Number of links dropped by the span due to limits being reached.
DroppedLinks uint32 
}

//go:wasmimport wasi:otel/tracing@0.2.0-draft on-start
func wasm_import_on_start(arg0 uintptr, arg1 uint32, arg2 uintptr, arg3 uint32, arg4 int32, arg5 int32, arg6 uintptr, arg7 uint32) 

func OnStart(context SpanContext)  {
        pinner := &runtime.Pinner{}
defer pinner.Unpin()

        
        utf8 := unsafe.Pointer(unsafe.StringData((context).TraceId))
pinner.Pin(utf8)
utf80 := unsafe.Pointer(unsafe.StringData((context).SpanId))
pinner.Pin(utf80)
var result int32
if (context).IsRemote {
        result = 1
} else {
        result = 0
}
slice := (context).TraceState
length := uint32(len(slice))
result3 := wit_runtime.Allocate(pinner, uintptr(length * (4*4)), 4)
for index, element := range slice {
        base := unsafe.Add(result3, index * (4*4))
        utf81 := unsafe.Pointer(unsafe.StringData((element).F0))
pinner.Pin(utf81)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).F0)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf81)))
utf82 := unsafe.Pointer(unsafe.StringData((element).F1))
pinner.Pin(utf82)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (3*4))) = uint32(uint32(len((element).F1)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (2*4))) = uint32(uintptr(uintptr(utf82)))

}

wasm_import_on_start(uintptr(utf8), uint32(len((context).TraceId)), uintptr(utf80), uint32(len((context).SpanId)), int32((context).TraceFlags), result, uintptr(result3), length)

}

//go:wasmimport wasi:otel/tracing@0.2.0-draft on-end
func wasm_import_on_end(arg0 uintptr) 

func OnEnd(span SpanData)  {
        pinner := &runtime.Pinner{}
defer pinner.Unpin()

        returnArea := uintptr(wit_runtime.Allocate(pinner, (40+32*4), 8))
        utf8 := unsafe.Pointer(unsafe.StringData(((span).SpanContext).TraceId))
pinner.Pin(utf8)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), 4)) = uint32(uint32(len(((span).SpanContext).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), 0)) = uint32(uintptr(uintptr(utf8)))
utf80 := unsafe.Pointer(unsafe.StringData(((span).SpanContext).SpanId))
pinner.Pin(utf80)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (3*4))) = uint32(uint32(len(((span).SpanContext).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (2*4))) = uint32(uintptr(uintptr(utf80)))
*(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (4*4))) = int8(int32(((span).SpanContext).TraceFlags))
var result int32
if ((span).SpanContext).IsRemote {
        result = 1
} else {
        result = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (1+4*4))) = int8(result)
slice := ((span).SpanContext).TraceState
length := uint32(len(slice))
result3 := wit_runtime.Allocate(pinner, uintptr(length * (4*4)), 4)
for index, element := range slice {
        base := unsafe.Add(result3, index * (4*4))
        utf81 := unsafe.Pointer(unsafe.StringData((element).F0))
pinner.Pin(utf81)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).F0)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf81)))
utf82 := unsafe.Pointer(unsafe.StringData((element).F1))
pinner.Pin(utf82)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (3*4))) = uint32(uint32(len((element).F1)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (2*4))) = uint32(uintptr(uintptr(utf82)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (6*4))) = uint32(length)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (5*4))) = uint32(uintptr(uintptr(result3)))
utf84 := unsafe.Pointer(unsafe.StringData((span).ParentSpanId))
pinner.Pin(utf84)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (8*4))) = uint32(uint32(len((span).ParentSpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (7*4))) = uint32(uintptr(uintptr(utf84)))
*(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (9*4))) = int8(int32((span).SpanKind))
utf85 := unsafe.Pointer(unsafe.StringData((span).Name))
pinner.Pin(utf85)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (11*4))) = uint32(uint32(len((span).Name)))
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (10*4))) = uint32(uintptr(uintptr(utf85)))
*(*int64)(unsafe.Add(unsafe.Pointer(returnArea), (12*4))) = int64(((span).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(returnArea), (8+12*4))) = int32(((span).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(returnArea), (16+12*4))) = int64(((span).EndTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(returnArea), (24+12*4))) = int32(((span).EndTime).Nanoseconds)
slice18 := (span).Attributes
length20 := uint32(len(slice18))
result19 := wit_runtime.Allocate(pinner, uintptr(length20 * (8+4*4)), 8)
for index, element := range slice18 {
        base := unsafe.Add(result19, index * (8+4*4))
        utf86 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf86)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf86)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf87 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf87)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf87)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result8 int32
if payload {
        result8 = 1
} else {
        result8 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result8)

        

case wasi_otel_types.ValueF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case wasi_otel_types.ValueS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(3))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case wasi_otel_types.ValueStringArray:
        payload := (element).Value.StringArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(4))
slice10 := payload
length12 := uint32(len(slice10))
result11 := wit_runtime.Allocate(pinner, uintptr(length12 * (2*4)), 4)
for index, element := range slice10 {
        base := unsafe.Add(result11, index * (2*4))
        utf89 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf89)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf89)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length12)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result11)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice14 := payload
length16 := uint32(len(slice14))
result15 := wit_runtime.Allocate(pinner, uintptr(length16 * 1), 1)
for index, element := range slice14 {
        base := unsafe.Add(result15, index * 1)
        var result13 int32
if element {
        result13 = 1
} else {
        result13 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result13)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length16)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result15)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data17 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data17)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data17)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+13*4))) = uint32(length20)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+12*4))) = uint32(uintptr(uintptr(result19)))
slice38 := (span).Events
length40 := uint32(len(slice38))
result39 := wit_runtime.Allocate(pinner, uintptr(length40 * (16+4*4)), 8)
for index, element := range slice38 {
        base := unsafe.Add(result39, index * (16+4*4))
        utf821 := unsafe.Pointer(unsafe.StringData((element).Name))
pinner.Pin(utf821)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Name)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf821)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)
slice35 := (element).Attributes
length37 := uint32(len(slice35))
result36 := wit_runtime.Allocate(pinner, uintptr(length37 * (8+4*4)), 8)
for index, element := range slice35 {
        base := unsafe.Add(result36, index * (8+4*4))
        utf822 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf822)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf822)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf823 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf823)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf823)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result24 int32
if payload {
        result24 = 1
} else {
        result24 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result24)

        

case wasi_otel_types.ValueF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case wasi_otel_types.ValueS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(3))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case wasi_otel_types.ValueStringArray:
        payload := (element).Value.StringArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(4))
slice26 := payload
length28 := uint32(len(slice26))
result27 := wit_runtime.Allocate(pinner, uintptr(length28 * (2*4)), 4)
for index, element := range slice26 {
        base := unsafe.Add(result27, index * (2*4))
        utf825 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf825)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf825)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length28)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result27)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice30 := payload
length32 := uint32(len(slice30))
result31 := wit_runtime.Allocate(pinner, uintptr(length32 * 1), 1)
for index, element := range slice30 {
        base := unsafe.Add(result31, index * 1)
        var result29 int32
if element {
        result29 = 1
} else {
        result29 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result29)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length32)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result31)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data33 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data33)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data33)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data34 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data34)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data34)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+3*4))) = uint32(length37)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = uint32(uintptr(uintptr(result36)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+15*4))) = uint32(length40)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+14*4))) = uint32(uintptr(uintptr(result39)))
slice65 := (span).Links
length67 := uint32(len(slice65))
result66 := wit_runtime.Allocate(pinner, uintptr(length67 * (9*4)), 4)
for index, element := range slice65 {
        base := unsafe.Add(result66, index * (9*4))
        utf841 := unsafe.Pointer(unsafe.StringData(((element).SpanContext).TraceId))
pinner.Pin(utf841)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(((element).SpanContext).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf841)))
utf842 := unsafe.Pointer(unsafe.StringData(((element).SpanContext).SpanId))
pinner.Pin(utf842)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (3*4))) = uint32(uint32(len(((element).SpanContext).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (2*4))) = uint32(uintptr(uintptr(utf842)))
*(*int8)(unsafe.Add(unsafe.Pointer(base), (4*4))) = int8(int32(((element).SpanContext).TraceFlags))
var result43 int32
if ((element).SpanContext).IsRemote {
        result43 = 1
} else {
        result43 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (1+4*4))) = int8(result43)
slice46 := ((element).SpanContext).TraceState
length48 := uint32(len(slice46))
result47 := wit_runtime.Allocate(pinner, uintptr(length48 * (4*4)), 4)
for index, element := range slice46 {
        base := unsafe.Add(result47, index * (4*4))
        utf844 := unsafe.Pointer(unsafe.StringData((element).F0))
pinner.Pin(utf844)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).F0)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf844)))
utf845 := unsafe.Pointer(unsafe.StringData((element).F1))
pinner.Pin(utf845)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (3*4))) = uint32(uint32(len((element).F1)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (2*4))) = uint32(uintptr(uintptr(utf845)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (6*4))) = uint32(length48)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (5*4))) = uint32(uintptr(uintptr(result47)))
slice62 := (element).Attributes
length64 := uint32(len(slice62))
result63 := wit_runtime.Allocate(pinner, uintptr(length64 * (8+4*4)), 8)
for index, element := range slice62 {
        base := unsafe.Add(result63, index * (8+4*4))
        utf849 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf849)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf849)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf850 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf850)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf850)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result51 int32
if payload {
        result51 = 1
} else {
        result51 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result51)

        

case wasi_otel_types.ValueF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case wasi_otel_types.ValueS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(3))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case wasi_otel_types.ValueStringArray:
        payload := (element).Value.StringArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(4))
slice53 := payload
length55 := uint32(len(slice53))
result54 := wit_runtime.Allocate(pinner, uintptr(length55 * (2*4)), 4)
for index, element := range slice53 {
        base := unsafe.Add(result54, index * (2*4))
        utf852 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf852)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf852)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length55)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result54)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice57 := payload
length59 := uint32(len(slice57))
result58 := wit_runtime.Allocate(pinner, uintptr(length59 * 1), 1)
for index, element := range slice57 {
        base := unsafe.Add(result58, index * 1)
        var result56 int32
if element {
        result56 = 1
} else {
        result56 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result56)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length59)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result58)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data60 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data60)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data60)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data61 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data61)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data61)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8*4))) = uint32(length64)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (7*4))) = uint32(uintptr(uintptr(result63)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+17*4))) = uint32(length67)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+16*4))) = uint32(uintptr(uintptr(result66)))

switch (span).Status.Tag() {
case StatusUnset:
        
        *(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (32+18*4))) = int8(int32(0))

        

case StatusOk:
        
        *(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (32+18*4))) = int8(int32(1))

        

case StatusError:
        payload := (span).Status.Error()
        *(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (32+18*4))) = int8(int32(2))
utf868 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf868)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+20*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+19*4))) = uint32(uintptr(uintptr(utf868)))

        

default:
        panic("unreachable")
}
utf869 := unsafe.Pointer(unsafe.StringData(((span).InstrumentationScope).Name))
pinner.Pin(utf869)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+22*4))) = uint32(uint32(len(((span).InstrumentationScope).Name)))
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+21*4))) = uint32(uintptr(uintptr(utf869)))

switch ((span).InstrumentationScope).Version.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (32+23*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := ((span).InstrumentationScope).Version.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (32+23*4))) = int8(int32(1))
utf870 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf870)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+25*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+24*4))) = uint32(uintptr(uintptr(utf870)))

        
default:
        panic("unreachable")
}

switch ((span).InstrumentationScope).SchemaUrl.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (32+26*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := ((span).InstrumentationScope).SchemaUrl.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(returnArea), (32+26*4))) = int8(int32(1))
utf871 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf871)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+28*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+27*4))) = uint32(uintptr(uintptr(utf871)))

        
default:
        panic("unreachable")
}
slice85 := ((span).InstrumentationScope).Attributes
length87 := uint32(len(slice85))
result86 := wit_runtime.Allocate(pinner, uintptr(length87 * (8+4*4)), 8)
for index, element := range slice85 {
        base := unsafe.Add(result86, index * (8+4*4))
        utf872 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf872)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf872)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf873 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf873)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf873)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result74 int32
if payload {
        result74 = 1
} else {
        result74 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result74)

        

case wasi_otel_types.ValueF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case wasi_otel_types.ValueS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(3))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case wasi_otel_types.ValueStringArray:
        payload := (element).Value.StringArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(4))
slice76 := payload
length78 := uint32(len(slice76))
result77 := wit_runtime.Allocate(pinner, uintptr(length78 * (2*4)), 4)
for index, element := range slice76 {
        base := unsafe.Add(result77, index * (2*4))
        utf875 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf875)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf875)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length78)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result77)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice80 := payload
length82 := uint32(len(slice80))
result81 := wit_runtime.Allocate(pinner, uintptr(length82 * 1), 1)
for index, element := range slice80 {
        base := unsafe.Add(result81, index * 1)
        var result79 int32
if element {
        result79 = 1
} else {
        result79 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result79)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length82)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result81)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data83 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data83)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data83)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data84 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data84)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data84)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+30*4))) = uint32(length87)
*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (32+29*4))) = uint32(uintptr(uintptr(result86)))
*(*int32)(unsafe.Add(unsafe.Pointer(returnArea), (32+31*4))) = int32((span).DroppedAttributes)
*(*int32)(unsafe.Add(unsafe.Pointer(returnArea), (36+31*4))) = int32((span).DroppedEvents)
*(*int32)(unsafe.Add(unsafe.Pointer(returnArea), (40+31*4))) = int32((span).DroppedLinks)
wasm_import_on_end(returnArea)

}

//go:wasmimport wasi:otel/tracing@0.2.0-draft outer-span-context
func wasm_import_outer_span_context(arg0 uintptr) 

func OuterSpanContext() SpanContext {
        pinner := &runtime.Pinner{}
defer pinner.Unpin()

        returnArea := uintptr(wit_runtime.Allocate(pinner, (7*4), 4))
        wasm_import_outer_span_context(returnArea)
value := unsafe.String((*uint8)(unsafe.Pointer(uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), 0))))), *(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), 4)))
value0 := unsafe.String((*uint8)(unsafe.Pointer(uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (2*4)))))), *(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (3*4))))
result := make([]wit_types.Tuple2[string, string], 0, *(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (6*4))))
for index := 0; index < int(*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (6*4)))); index++ {
        base := unsafe.Add(unsafe.Pointer(uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (5*4))))), index * (4*4))
        value1 := unsafe.String((*uint8)(unsafe.Pointer(uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0))))), *(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)))
value2 := unsafe.String((*uint8)(unsafe.Pointer(uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(base), (2*4)))))), *(*uint32)(unsafe.Add(unsafe.Pointer(base), (3*4))))

        result = append(result, wit_types.Tuple2[string, string]{value1, value2})        
}

result3 := SpanContext{value, value0, uint8(uint8(*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (4*4))))), (uint8(*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (1+4*4)))) != 0), result};
return result3

}

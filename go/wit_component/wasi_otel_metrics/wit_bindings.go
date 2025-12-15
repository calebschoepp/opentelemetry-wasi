package wasi_otel_metrics

import (
        "wit_component/wasi_clocks_wall_clock"
"wit_component/wasi_otel_types"
"wit_component/wit_runtime"
"wit_component/wit_types"
"unsafe"
"runtime"
)

type Datetime = wasi_clocks_wall_clock.Datetime
type Duration = uint64
type KeyValue = wasi_otel_types.KeyValue
type InstrumentationScope = wasi_otel_types.InstrumentationScope
type SpanId = string
type TraceId = string
// An error resulting from `export` being called.
type Error = string

// An immutable representation of the entity producing telemetry as attributes.
type Resource struct {
        // Attributes that identify the resource.
Attributes []wasi_otel_types.KeyValue
// The schema URL to be recorded in the emitted resource.
SchemaUrl wit_types.Option[string] 
}

// A set of bucket counts, encoded in a contiguous array of counts.
type ExponentialBucket struct {
        // The bucket index of the first entry in the `counts` list.
Offset int32
// A list where `counts[i]` carries the count of the bucket at index `offset + i`.
// 
// `counts[i]` is the count of values greater than base^(offset+i) and less than
// or equal to base^(offset+i+1).
Counts []uint64 
}

const (
        // A measurement interval that continues to expand forward in time from a
// starting point.
// 
// New measurements are added to all previous measurements since a start time.
// 
// This is the default temporality.
TemporalityCumulative uint8 = 0
// A measurement interval that resets each cycle.
// 
// Measurements from one cycle are recorded independently, measurements from
// other cycles do not affect them.
TemporalityDelta uint8 = 1
// Configures Synchronous Counter and Histogram instruments to use
// Delta aggregation temporality, which allows them to shed memory
// following a cardinality explosion, thus use less memory.
TemporalityLowMemory uint8 = 2
)
// Defines the window that an aggregation was calculated over.
type Temporality = uint8

const (
MetricNumberF64 uint8 = 0
MetricNumberS64 uint8 = 1
MetricNumberU64 uint8 = 2
)

// The number types available for any given instrument.
type MetricNumber struct {
        tag uint8
        value any
}

func (self MetricNumber) Tag() uint8 {
        return self.tag
}

func (self MetricNumber) F64() float64 {
        if self.tag != MetricNumberF64 {
                panic("tag mismatch")
        }
        return self.value.(float64)
}
func (self MetricNumber) S64() int64 {
        if self.tag != MetricNumberS64 {
                panic("tag mismatch")
        }
        return self.value.(int64)
}
func (self MetricNumber) U64() uint64 {
        if self.tag != MetricNumberU64 {
                panic("tag mismatch")
        }
        return self.value.(uint64)
}

func MakeMetricNumberF64(value float64) MetricNumber {
        return MetricNumber{MetricNumberF64, value}
}
func MakeMetricNumberS64(value int64) MetricNumber {
        return MetricNumber{MetricNumberS64, value}
}
func MakeMetricNumberU64(value uint64) MetricNumber {
        return MetricNumber{MetricNumberU64, value}
}



// A measurement sampled from a time series providing a typical example.
type Exemplar struct {
        // The attributes recorded with the measurement but filtered out of the
// time series' aggregated data.
FilteredAttributes []wasi_otel_types.KeyValue
// The time when the measurement was recorded.
Time wasi_clocks_wall_clock.Datetime
// The measured value.
Value MetricNumber
// The ID of the span that was active during the measurement.
// 
// If no span was active or the span was not sampled this will be empty.
SpanId string
// The ID of the trace the active span belonged to during the measurement.
// 
// If no span was active or the span was not sampled this will be empty.
TraceId string 
}

// A single data point in a time series to be associated with a `gauge`.
type GaugeDataPoint struct {
        // `attributes` is the set of key value pairs that uniquely identify the
// time series.
Attributes []wasi_otel_types.KeyValue
// The value of this data point.
Value MetricNumber
// The sampled `exemplar`s collected during the time series.
Exemplars []Exemplar 
}

// A measurement of the current value of an instrument.
type Gauge struct {
        // Represents individual aggregated measurements with unique attributes.
DataPoints []GaugeDataPoint
// The time when the time series was started.
StartTime wit_types.Option[wasi_clocks_wall_clock.Datetime]
// The time when the time series was recorded.
Time wasi_clocks_wall_clock.Datetime 
}

// A single data point in a time series to be associated with a `sum`.
type SumDataPoint struct {
        // `attributes` is the set of key value pairs that uniquely identify the
// time series.
Attributes []wasi_otel_types.KeyValue
// The value of this data point.
Value MetricNumber
// The sampled `exemplar`s collected during the time series.
Exemplars []Exemplar 
}

// Represents the sum of all measurements of values from an instrument.
type Sum struct {
        // Represents individual aggregated measurements with unique attributes.
DataPoints []SumDataPoint
// The time when the time series was started.
StartTime wasi_clocks_wall_clock.Datetime
// The time when the time series was recorded.
Time wasi_clocks_wall_clock.Datetime
// Describes if the aggregation is reported as the change from the last report
// time, or the cumulative changes since a fixed start time.
Temporality Temporality
// Whether this aggregation only increases or decreases.
IsMonotonic bool 
}

// A single data point in a time series to be associated with a `histogram`.
type HistogramDataPoint struct {
        // The set of key value pairs that uniquely identify the time series.
Attributes []wasi_otel_types.KeyValue
// The number of updates this histogram has been calculated with.
Count uint64
// The upper bounds of the buckets of the histogram.
Bounds []float64
// The count of each of the buckets.
BucketCounts []uint64
// The minimum value recorded.
Min wit_types.Option[MetricNumber]
// The maximum value recorded.
Max wit_types.Option[MetricNumber]
// The sum of the values recorded
Sum MetricNumber
// The sampled `exemplar`s collected during the time series.
Exemplars []Exemplar 
}

// Represents the histogram of all measurements of values from an instrument.
type Histogram struct {
        // Individual aggregated measurements with unique attributes.
DataPoints []HistogramDataPoint
// The time when the time series was started.
StartTime wasi_clocks_wall_clock.Datetime
// The time when the time series was recorded.
Time wasi_clocks_wall_clock.Datetime
// Describes if the aggregation is reported as the change from the last report
// time, or the cumulative changes since a fixed start time.
Temporality Temporality 
}

// A single data point in a time series to be associated with an `exponential-histogram `.
type ExponentialHistogramDataPoint struct {
        // The set of key value pairs that uniquely identify the time series.
Attributes []wasi_otel_types.KeyValue
// The number of updates this histogram has been calculated with.
Count uint64
// The minimum value recorded.
Min wit_types.Option[MetricNumber]
// The maximum value recorded.
Max wit_types.Option[MetricNumber]
// The maximum value recorded.
Sum MetricNumber
// Describes the resolution of the histogram.
// 
// Boundaries are located at powers of the base, where:
// 
//   base = 2 ^ (2 ^ -scale)
Scale int8
// The number of values whose absolute value is less than or equal to
// `zero_threshold`.
// 
// When `zero_threshold` is `0`, this is the number of values that cannot be
// expressed using the standard exponential formula as well as values that have
// been rounded to zero.
ZeroCount uint64
// The range of positive value bucket counts.
PositiveBucket ExponentialBucket
// The range of negative value bucket counts.
NegativeBucket ExponentialBucket
// The width of the zero region.
// 
// Where the zero region is defined as the closed interval
// [-zero_threshold, zero_threshold].
ZeroThreshold float64
// The sampled exemplars collected during the time series.
Exemplars []Exemplar 
}

// The histogram of all measurements of values from an instrument.
type ExponentialHistogram struct {
        // The individual aggregated measurements with unique attributes.
DataPoints []ExponentialHistogramDataPoint
// When the time series was started.
StartTime wasi_clocks_wall_clock.Datetime
// The time when the time series was recorded.
Time wasi_clocks_wall_clock.Datetime
// Describes if the aggregation is reported as the change from the last report
// time, or the cumulative changes since a fixed start time.
Temporality Temporality 
}

const (
// Metric data for an f64 gauge.
MetricDataF64Gauge uint8 = 0
// Metric data for an f64 sum.
MetricDataF64Sum uint8 = 1
// Metric data for an f64 histogram.
MetricDataF64Histogram uint8 = 2
// Metric data for an f64 exponential-histogram.
MetricDataF64ExponentialHistogram uint8 = 3
// Metric data for an u64 gauge.
MetricDataU64Gauge uint8 = 4
// Metric data for an u64 sum.
MetricDataU64Sum uint8 = 5
// Metric data for an u64 histogram.
MetricDataU64Histogram uint8 = 6
// Metric data for an u64 exponential-histogram.
MetricDataU64ExponentialHistogram uint8 = 7
// Metric data for an s64 gauge.
MetricDataS64Gauge uint8 = 8
// Metric data for an s64 sum.
MetricDataS64Sum uint8 = 9
// Metric data for an s64 histogram.
MetricDataS64Histogram uint8 = 10
// Metric data for an s64 exponential-histogram.
MetricDataS64ExponentialHistogram uint8 = 11
)

// Metric data for all types.
type MetricData struct {
        tag uint8
        value any
}

func (self MetricData) Tag() uint8 {
        return self.tag
}

func (self MetricData) F64Gauge() Gauge {
        if self.tag != MetricDataF64Gauge {
                panic("tag mismatch")
        }
        return self.value.(Gauge)
}
func (self MetricData) F64Sum() Sum {
        if self.tag != MetricDataF64Sum {
                panic("tag mismatch")
        }
        return self.value.(Sum)
}
func (self MetricData) F64Histogram() Histogram {
        if self.tag != MetricDataF64Histogram {
                panic("tag mismatch")
        }
        return self.value.(Histogram)
}
func (self MetricData) F64ExponentialHistogram() ExponentialHistogram {
        if self.tag != MetricDataF64ExponentialHistogram {
                panic("tag mismatch")
        }
        return self.value.(ExponentialHistogram)
}
func (self MetricData) U64Gauge() Gauge {
        if self.tag != MetricDataU64Gauge {
                panic("tag mismatch")
        }
        return self.value.(Gauge)
}
func (self MetricData) U64Sum() Sum {
        if self.tag != MetricDataU64Sum {
                panic("tag mismatch")
        }
        return self.value.(Sum)
}
func (self MetricData) U64Histogram() Histogram {
        if self.tag != MetricDataU64Histogram {
                panic("tag mismatch")
        }
        return self.value.(Histogram)
}
func (self MetricData) U64ExponentialHistogram() ExponentialHistogram {
        if self.tag != MetricDataU64ExponentialHistogram {
                panic("tag mismatch")
        }
        return self.value.(ExponentialHistogram)
}
func (self MetricData) S64Gauge() Gauge {
        if self.tag != MetricDataS64Gauge {
                panic("tag mismatch")
        }
        return self.value.(Gauge)
}
func (self MetricData) S64Sum() Sum {
        if self.tag != MetricDataS64Sum {
                panic("tag mismatch")
        }
        return self.value.(Sum)
}
func (self MetricData) S64Histogram() Histogram {
        if self.tag != MetricDataS64Histogram {
                panic("tag mismatch")
        }
        return self.value.(Histogram)
}
func (self MetricData) S64ExponentialHistogram() ExponentialHistogram {
        if self.tag != MetricDataS64ExponentialHistogram {
                panic("tag mismatch")
        }
        return self.value.(ExponentialHistogram)
}

func MakeMetricDataF64Gauge(value Gauge) MetricData {
        return MetricData{MetricDataF64Gauge, value}
}
func MakeMetricDataF64Sum(value Sum) MetricData {
        return MetricData{MetricDataF64Sum, value}
}
func MakeMetricDataF64Histogram(value Histogram) MetricData {
        return MetricData{MetricDataF64Histogram, value}
}
func MakeMetricDataF64ExponentialHistogram(value ExponentialHistogram) MetricData {
        return MetricData{MetricDataF64ExponentialHistogram, value}
}
func MakeMetricDataU64Gauge(value Gauge) MetricData {
        return MetricData{MetricDataU64Gauge, value}
}
func MakeMetricDataU64Sum(value Sum) MetricData {
        return MetricData{MetricDataU64Sum, value}
}
func MakeMetricDataU64Histogram(value Histogram) MetricData {
        return MetricData{MetricDataU64Histogram, value}
}
func MakeMetricDataU64ExponentialHistogram(value ExponentialHistogram) MetricData {
        return MetricData{MetricDataU64ExponentialHistogram, value}
}
func MakeMetricDataS64Gauge(value Gauge) MetricData {
        return MetricData{MetricDataS64Gauge, value}
}
func MakeMetricDataS64Sum(value Sum) MetricData {
        return MetricData{MetricDataS64Sum, value}
}
func MakeMetricDataS64Histogram(value Histogram) MetricData {
        return MetricData{MetricDataS64Histogram, value}
}
func MakeMetricDataS64ExponentialHistogram(value ExponentialHistogram) MetricData {
        return MetricData{MetricDataS64ExponentialHistogram, value}
}



// A collection of one or more aggregated time series from a metric.
type Metric struct {
        // The name of the metric that created this data.
Name string
// The description of the metric, which can be used in documentation.
Description string
// The unit in which the metric reports.
Unit string
// The aggregated data from a metric.
Data MetricData 
}

// A collection of `metric`s produced by a meter.
type ScopeMetrics struct {
        // The instrumentation scope that the meter was created with.
Scope wasi_otel_types.InstrumentationScope
// The list of aggregations created by the meter.
Metrics []Metric 
}

// A collection of `scope-metrics` and the associated `resource` that created them.
type ResourceMetrics struct {
        // The entity that collected the metrics.
Resource Resource
// The collection of metrics with unique `instrumentation-scope`s.
ScopeMetrics []ScopeMetrics 
}

//go:wasmimport wasi:otel/metrics@0.2.0-draft export
func wasm_import_export(arg0 uintptr, arg1 uint32, arg2 int32, arg3 uintptr, arg4 uint32, arg5 uintptr, arg6 uint32, arg7 uintptr) 

func Export(metrics ResourceMetrics) wit_types.Result[wit_types.Unit, string] {
        pinner := &runtime.Pinner{}
defer pinner.Unpin()

        returnArea := uintptr(wit_runtime.Allocate(pinner, (3*4), 4))
        slice8 := ((metrics).Resource).Attributes
length10 := uint32(len(slice8))
result9 := wit_runtime.Allocate(pinner, uintptr(length10 * (8+4*4)), 8)
for index, element := range slice8 {
        base := unsafe.Add(result9, index * (8+4*4))
        utf8 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf80 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf80)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf80)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result int32
if payload {
        result = 1
} else {
        result = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result)

        

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
slice := payload
length := uint32(len(slice))
result2 := wit_runtime.Allocate(pinner, uintptr(length * (2*4)), 4)
for index, element := range slice {
        base := unsafe.Add(result2, index * (2*4))
        utf81 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf81)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf81)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result2)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice4 := payload
length6 := uint32(len(slice4))
result5 := wit_runtime.Allocate(pinner, uintptr(length6 * 1), 1)
for index, element := range slice4 {
        base := unsafe.Add(result5, index * 1)
        var result3 int32
if element {
        result3 = 1
} else {
        result3 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result3)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length6)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result5)))

        

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
data7 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data7)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data7)))

        

default:
        panic("unreachable")
}

}

var option int32
var option12 uintptr
var option13 uint32
switch ((metrics).Resource).SchemaUrl.Tag() {
case wit_types.OptionNone:
        
        option = int32(0)
option12 = 0
option13 = 0
case wit_types.OptionSome:
        payload := ((metrics).Resource).SchemaUrl.Some()
        utf811 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf811)

        option = int32(1)
option12 = uintptr(utf811)
option13 = uint32(len(payload))
default:
        panic("unreachable")
}
slice534 := (metrics).ScopeMetrics
length536 := uint32(len(slice534))
result535 := wit_runtime.Allocate(pinner, uintptr(length536 * (12*4)), 4)
for index, element := range slice534 {
        base := unsafe.Add(result535, index * (12*4))
        utf814 := unsafe.Pointer(unsafe.StringData(((element).Scope).Name))
pinner.Pin(utf814)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(((element).Scope).Name)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf814)))

switch ((element).Scope).Version.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := ((element).Scope).Version.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
utf815 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf815)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (4*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (3*4))) = uint32(uintptr(uintptr(utf815)))

        
default:
        panic("unreachable")
}

switch ((element).Scope).SchemaUrl.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (5*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := ((element).Scope).SchemaUrl.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (5*4))) = int8(int32(1))
utf816 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf816)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (7*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (6*4))) = uint32(uintptr(uintptr(utf816)))

        
default:
        panic("unreachable")
}
slice30 := ((element).Scope).Attributes
length32 := uint32(len(slice30))
result31 := wit_runtime.Allocate(pinner, uintptr(length32 * (8+4*4)), 8)
for index, element := range slice30 {
        base := unsafe.Add(result31, index * (8+4*4))
        utf817 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf817)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf817)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf818 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf818)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf818)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result19 int32
if payload {
        result19 = 1
} else {
        result19 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result19)

        

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
slice21 := payload
length23 := uint32(len(slice21))
result22 := wit_runtime.Allocate(pinner, uintptr(length23 * (2*4)), 4)
for index, element := range slice21 {
        base := unsafe.Add(result22, index * (2*4))
        utf820 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf820)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf820)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length23)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result22)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice25 := payload
length27 := uint32(len(slice25))
result26 := wit_runtime.Allocate(pinner, uintptr(length27 * 1), 1)
for index, element := range slice25 {
        base := unsafe.Add(result26, index * 1)
        var result24 int32
if element {
        result24 = 1
} else {
        result24 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result24)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length27)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result26)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data28 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data28)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data28)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data29 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data29)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data29)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (9*4))) = uint32(length32)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8*4))) = uint32(uintptr(uintptr(result31)))
slice531 := (element).Metrics
length533 := uint32(len(slice531))
result532 := wit_runtime.Allocate(pinner, uintptr(length533 * (48+8*4)), 8)
for index, element := range slice531 {
        base := unsafe.Add(result532, index * (48+8*4))
        utf833 := unsafe.Pointer(unsafe.StringData((element).Name))
pinner.Pin(utf833)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Name)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf833)))
utf834 := unsafe.Pointer(unsafe.StringData((element).Description))
pinner.Pin(utf834)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (3*4))) = uint32(uint32(len((element).Description)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (2*4))) = uint32(uintptr(uintptr(utf834)))
utf835 := unsafe.Pointer(unsafe.StringData((element).Unit))
pinner.Pin(utf835)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (5*4))) = uint32(uint32(len((element).Unit)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (4*4))) = uint32(uintptr(uintptr(utf835)))

switch (element).Data.Tag() {
case MetricDataF64Gauge:
        payload := (element).Data.F64Gauge()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(0))
slice73 := (payload).DataPoints
length75 := uint32(len(slice73))
result74 := wit_runtime.Allocate(pinner, uintptr(length75 * (16+4*4)), 8)
for index, element := range slice73 {
        base := unsafe.Add(result74, index * (16+4*4))
        slice49 := (element).Attributes
length51 := uint32(len(slice49))
result50 := wit_runtime.Allocate(pinner, uintptr(length51 * (8+4*4)), 8)
for index, element := range slice49 {
        base := unsafe.Add(result50, index * (8+4*4))
        utf836 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf836)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf836)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf837 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf837)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf837)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result38 int32
if payload {
        result38 = 1
} else {
        result38 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result38)

        

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
slice40 := payload
length42 := uint32(len(slice40))
result41 := wit_runtime.Allocate(pinner, uintptr(length42 * (2*4)), 4)
for index, element := range slice40 {
        base := unsafe.Add(result41, index * (2*4))
        utf839 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf839)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf839)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length42)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result41)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice44 := payload
length46 := uint32(len(slice44))
result45 := wit_runtime.Allocate(pinner, uintptr(length46 * 1), 1)
for index, element := range slice44 {
        base := unsafe.Add(result45, index * 1)
        var result43 int32
if element {
        result43 = 1
} else {
        result43 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result43)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length46)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result45)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data47 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data47)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data47)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data48 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data48)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data48)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length51)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result50)))

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice70 := (element).Exemplars
length72 := uint32(len(slice70))
result71 := wit_runtime.Allocate(pinner, uintptr(length72 * (32+6*4)), 8)
for index, element := range slice70 {
        base := unsafe.Add(result71, index * (32+6*4))
        slice65 := (element).FilteredAttributes
length67 := uint32(len(slice65))
result66 := wit_runtime.Allocate(pinner, uintptr(length67 * (8+4*4)), 8)
for index, element := range slice65 {
        base := unsafe.Add(result66, index * (8+4*4))
        utf852 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf852)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf852)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf853 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf853)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf853)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result54 int32
if payload {
        result54 = 1
} else {
        result54 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result54)

        

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
slice56 := payload
length58 := uint32(len(slice56))
result57 := wit_runtime.Allocate(pinner, uintptr(length58 * (2*4)), 4)
for index, element := range slice56 {
        base := unsafe.Add(result57, index * (2*4))
        utf855 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf855)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf855)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length58)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result57)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice60 := payload
length62 := uint32(len(slice60))
result61 := wit_runtime.Allocate(pinner, uintptr(length62 * 1), 1)
for index, element := range slice60 {
        base := unsafe.Add(result61, index * 1)
        var result59 int32
if element {
        result59 = 1
} else {
        result59 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result59)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length62)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result61)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data63 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data63)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data63)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data64 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data64)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data64)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length67)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result66)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf868 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf868)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf868)))
utf869 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf869)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf869)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+3*4))) = uint32(length72)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = uint32(uintptr(uintptr(result71)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length75)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result74)))

switch (payload).StartTime.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (payload).StartTime.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int64((payload).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int32((payload).Nanoseconds)

        
default:
        panic("unreachable")
}
*(*int64)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int32(((payload).Time).Nanoseconds)

        

case MetricDataF64Sum:
        payload := (element).Data.F64Sum()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(1))
slice113 := (payload).DataPoints
length115 := uint32(len(slice113))
result114 := wit_runtime.Allocate(pinner, uintptr(length115 * (16+4*4)), 8)
for index, element := range slice113 {
        base := unsafe.Add(result114, index * (16+4*4))
        slice89 := (element).Attributes
length91 := uint32(len(slice89))
result90 := wit_runtime.Allocate(pinner, uintptr(length91 * (8+4*4)), 8)
for index, element := range slice89 {
        base := unsafe.Add(result90, index * (8+4*4))
        utf876 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf876)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf876)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf877 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf877)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf877)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result78 int32
if payload {
        result78 = 1
} else {
        result78 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result78)

        

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
slice80 := payload
length82 := uint32(len(slice80))
result81 := wit_runtime.Allocate(pinner, uintptr(length82 * (2*4)), 4)
for index, element := range slice80 {
        base := unsafe.Add(result81, index * (2*4))
        utf879 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf879)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf879)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length82)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result81)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice84 := payload
length86 := uint32(len(slice84))
result85 := wit_runtime.Allocate(pinner, uintptr(length86 * 1), 1)
for index, element := range slice84 {
        base := unsafe.Add(result85, index * 1)
        var result83 int32
if element {
        result83 = 1
} else {
        result83 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result83)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length86)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result85)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data87 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data87)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data87)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data88 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data88)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data88)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length91)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result90)))

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice110 := (element).Exemplars
length112 := uint32(len(slice110))
result111 := wit_runtime.Allocate(pinner, uintptr(length112 * (32+6*4)), 8)
for index, element := range slice110 {
        base := unsafe.Add(result111, index * (32+6*4))
        slice105 := (element).FilteredAttributes
length107 := uint32(len(slice105))
result106 := wit_runtime.Allocate(pinner, uintptr(length107 * (8+4*4)), 8)
for index, element := range slice105 {
        base := unsafe.Add(result106, index * (8+4*4))
        utf892 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf892)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf892)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf893 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf893)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf893)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result94 int32
if payload {
        result94 = 1
} else {
        result94 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result94)

        

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
slice96 := payload
length98 := uint32(len(slice96))
result97 := wit_runtime.Allocate(pinner, uintptr(length98 * (2*4)), 4)
for index, element := range slice96 {
        base := unsafe.Add(result97, index * (2*4))
        utf895 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf895)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf895)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length98)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result97)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice100 := payload
length102 := uint32(len(slice100))
result101 := wit_runtime.Allocate(pinner, uintptr(length102 * 1), 1)
for index, element := range slice100 {
        base := unsafe.Add(result101, index * 1)
        var result99 int32
if element {
        result99 = 1
} else {
        result99 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result99)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length102)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result101)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data103 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data103)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data103)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data104 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data104)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data104)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length107)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result106)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8108 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8108)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8108)))
utf8109 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8109)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8109)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+3*4))) = uint32(length112)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = uint32(uintptr(uintptr(result111)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length115)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result114)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))
var result116 int32
if (payload).IsMonotonic {
        result116 = 1
} else {
        result116 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (41+8*4))) = int8(result116)

        

case MetricDataF64Histogram:
        payload := (element).Data.F64Histogram()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(2))
slice156 := (payload).DataPoints
length158 := uint32(len(slice156))
result157 := wit_runtime.Allocate(pinner, uintptr(length158 * (72+8*4)), 8)
for index, element := range slice156 {
        base := unsafe.Add(result157, index * (72+8*4))
        slice130 := (element).Attributes
length132 := uint32(len(slice130))
result131 := wit_runtime.Allocate(pinner, uintptr(length132 * (8+4*4)), 8)
for index, element := range slice130 {
        base := unsafe.Add(result131, index * (8+4*4))
        utf8117 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8117)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8117)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8118 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8118)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8118)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result119 int32
if payload {
        result119 = 1
} else {
        result119 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result119)

        

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
slice121 := payload
length123 := uint32(len(slice121))
result122 := wit_runtime.Allocate(pinner, uintptr(length123 * (2*4)), 4)
for index, element := range slice121 {
        base := unsafe.Add(result122, index * (2*4))
        utf8120 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8120)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8120)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length123)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result122)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice125 := payload
length127 := uint32(len(slice125))
result126 := wit_runtime.Allocate(pinner, uintptr(length127 * 1), 1)
for index, element := range slice125 {
        base := unsafe.Add(result126, index * 1)
        var result124 int32
if element {
        result124 = 1
} else {
        result124 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result124)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length127)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result126)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data128 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data128)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data128)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data129 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data129)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data129)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length132)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result131)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64((element).Count)
data133 := unsafe.Pointer(unsafe.SliceData((element).Bounds))
pinner.Pin(data133)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len((element).Bounds)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data133)))
data134 := unsafe.Pointer(unsafe.SliceData((element).BucketCounts))
pinner.Pin(data134)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+5*4))) = uint32(uint32(len((element).BucketCounts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+4*4))) = uint32(uintptr(uintptr(data134)))

switch (element).Min.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Min.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Max.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+6*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Max.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+6*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Sum.Tag() {
case MetricNumberF64:
        payload := (element).Sum.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = payload

        

case MetricNumberS64:
        payload := (element).Sum.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = payload

        

case MetricNumberU64:
        payload := (element).Sum.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice153 := (element).Exemplars
length155 := uint32(len(slice153))
result154 := wit_runtime.Allocate(pinner, uintptr(length155 * (32+6*4)), 8)
for index, element := range slice153 {
        base := unsafe.Add(result154, index * (32+6*4))
        slice148 := (element).FilteredAttributes
length150 := uint32(len(slice148))
result149 := wit_runtime.Allocate(pinner, uintptr(length150 * (8+4*4)), 8)
for index, element := range slice148 {
        base := unsafe.Add(result149, index * (8+4*4))
        utf8135 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8135)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8135)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8136 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8136)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8136)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result137 int32
if payload {
        result137 = 1
} else {
        result137 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result137)

        

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
slice139 := payload
length141 := uint32(len(slice139))
result140 := wit_runtime.Allocate(pinner, uintptr(length141 * (2*4)), 4)
for index, element := range slice139 {
        base := unsafe.Add(result140, index * (2*4))
        utf8138 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8138)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8138)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length141)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result140)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice143 := payload
length145 := uint32(len(slice143))
result144 := wit_runtime.Allocate(pinner, uintptr(length145 * 1), 1)
for index, element := range slice143 {
        base := unsafe.Add(result144, index * 1)
        var result142 int32
if element {
        result142 = 1
} else {
        result142 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result142)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length145)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result144)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data146 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data146)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data146)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data147 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data147)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data147)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length150)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result149)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8151 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8151)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8151)))
utf8152 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8152)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8152)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (72+7*4))) = uint32(length155)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (72+6*4))) = uint32(uintptr(uintptr(result154)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length158)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result157)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))

        

case MetricDataF64ExponentialHistogram:
        payload := (element).Data.F64ExponentialHistogram()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(3))
slice198 := (payload).DataPoints
length200 := uint32(len(slice198))
result199 := wit_runtime.Allocate(pinner, uintptr(length200 * (96+10*4)), 8)
for index, element := range slice198 {
        base := unsafe.Add(result199, index * (96+10*4))
        slice172 := (element).Attributes
length174 := uint32(len(slice172))
result173 := wit_runtime.Allocate(pinner, uintptr(length174 * (8+4*4)), 8)
for index, element := range slice172 {
        base := unsafe.Add(result173, index * (8+4*4))
        utf8159 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8159)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8159)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8160 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8160)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8160)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result161 int32
if payload {
        result161 = 1
} else {
        result161 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result161)

        

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
slice163 := payload
length165 := uint32(len(slice163))
result164 := wit_runtime.Allocate(pinner, uintptr(length165 * (2*4)), 4)
for index, element := range slice163 {
        base := unsafe.Add(result164, index * (2*4))
        utf8162 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8162)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8162)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length165)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result164)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice167 := payload
length169 := uint32(len(slice167))
result168 := wit_runtime.Allocate(pinner, uintptr(length169 * 1), 1)
for index, element := range slice167 {
        base := unsafe.Add(result168, index * 1)
        var result166 int32
if element {
        result166 = 1
} else {
        result166 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result166)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length169)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result168)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data170 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data170)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data170)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data171 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data171)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data171)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length174)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result173)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64((element).Count)

switch (element).Min.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Min.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Max.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Max.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Sum.Tag() {
case MetricNumberF64:
        payload := (element).Sum.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Sum.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Sum.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (72+2*4))) = int8(int32((element).Scale))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (80+2*4))) = int64((element).ZeroCount)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (88+2*4))) = ((element).PositiveBucket).Offset
data175 := unsafe.Pointer(unsafe.SliceData(((element).PositiveBucket).Counts))
pinner.Pin(data175)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+4*4))) = uint32(uint32(len(((element).PositiveBucket).Counts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+3*4))) = uint32(uintptr(uintptr(data175)))
*(*int32)(unsafe.Add(unsafe.Pointer(base), (88+5*4))) = ((element).NegativeBucket).Offset
data176 := unsafe.Pointer(unsafe.SliceData(((element).NegativeBucket).Counts))
pinner.Pin(data176)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+7*4))) = uint32(uint32(len(((element).NegativeBucket).Counts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+6*4))) = uint32(uintptr(uintptr(data176)))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (88+8*4))) = (element).ZeroThreshold
slice195 := (element).Exemplars
length197 := uint32(len(slice195))
result196 := wit_runtime.Allocate(pinner, uintptr(length197 * (32+6*4)), 8)
for index, element := range slice195 {
        base := unsafe.Add(result196, index * (32+6*4))
        slice190 := (element).FilteredAttributes
length192 := uint32(len(slice190))
result191 := wit_runtime.Allocate(pinner, uintptr(length192 * (8+4*4)), 8)
for index, element := range slice190 {
        base := unsafe.Add(result191, index * (8+4*4))
        utf8177 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8177)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8177)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8178 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8178)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8178)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result179 int32
if payload {
        result179 = 1
} else {
        result179 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result179)

        

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
slice181 := payload
length183 := uint32(len(slice181))
result182 := wit_runtime.Allocate(pinner, uintptr(length183 * (2*4)), 4)
for index, element := range slice181 {
        base := unsafe.Add(result182, index * (2*4))
        utf8180 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8180)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8180)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length183)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result182)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice185 := payload
length187 := uint32(len(slice185))
result186 := wit_runtime.Allocate(pinner, uintptr(length187 * 1), 1)
for index, element := range slice185 {
        base := unsafe.Add(result186, index * 1)
        var result184 int32
if element {
        result184 = 1
} else {
        result184 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result184)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length187)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result186)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data188 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data188)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data188)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data189 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data189)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data189)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length192)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result191)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8193 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8193)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8193)))
utf8194 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8194)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8194)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (96+9*4))) = uint32(length197)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (96+8*4))) = uint32(uintptr(uintptr(result196)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length200)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result199)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))

        

case MetricDataU64Gauge:
        payload := (element).Data.U64Gauge()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(4))
slice238 := (payload).DataPoints
length240 := uint32(len(slice238))
result239 := wit_runtime.Allocate(pinner, uintptr(length240 * (16+4*4)), 8)
for index, element := range slice238 {
        base := unsafe.Add(result239, index * (16+4*4))
        slice214 := (element).Attributes
length216 := uint32(len(slice214))
result215 := wit_runtime.Allocate(pinner, uintptr(length216 * (8+4*4)), 8)
for index, element := range slice214 {
        base := unsafe.Add(result215, index * (8+4*4))
        utf8201 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8201)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8201)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8202 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8202)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8202)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result203 int32
if payload {
        result203 = 1
} else {
        result203 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result203)

        

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
slice205 := payload
length207 := uint32(len(slice205))
result206 := wit_runtime.Allocate(pinner, uintptr(length207 * (2*4)), 4)
for index, element := range slice205 {
        base := unsafe.Add(result206, index * (2*4))
        utf8204 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8204)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8204)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length207)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result206)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice209 := payload
length211 := uint32(len(slice209))
result210 := wit_runtime.Allocate(pinner, uintptr(length211 * 1), 1)
for index, element := range slice209 {
        base := unsafe.Add(result210, index * 1)
        var result208 int32
if element {
        result208 = 1
} else {
        result208 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result208)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length211)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result210)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data212 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data212)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data212)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data213 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data213)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data213)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length216)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result215)))

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice235 := (element).Exemplars
length237 := uint32(len(slice235))
result236 := wit_runtime.Allocate(pinner, uintptr(length237 * (32+6*4)), 8)
for index, element := range slice235 {
        base := unsafe.Add(result236, index * (32+6*4))
        slice230 := (element).FilteredAttributes
length232 := uint32(len(slice230))
result231 := wit_runtime.Allocate(pinner, uintptr(length232 * (8+4*4)), 8)
for index, element := range slice230 {
        base := unsafe.Add(result231, index * (8+4*4))
        utf8217 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8217)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8217)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8218 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8218)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8218)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result219 int32
if payload {
        result219 = 1
} else {
        result219 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result219)

        

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
slice221 := payload
length223 := uint32(len(slice221))
result222 := wit_runtime.Allocate(pinner, uintptr(length223 * (2*4)), 4)
for index, element := range slice221 {
        base := unsafe.Add(result222, index * (2*4))
        utf8220 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8220)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8220)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length223)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result222)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice225 := payload
length227 := uint32(len(slice225))
result226 := wit_runtime.Allocate(pinner, uintptr(length227 * 1), 1)
for index, element := range slice225 {
        base := unsafe.Add(result226, index * 1)
        var result224 int32
if element {
        result224 = 1
} else {
        result224 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result224)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length227)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result226)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data228 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data228)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data228)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data229 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data229)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data229)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length232)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result231)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8233 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8233)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8233)))
utf8234 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8234)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8234)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+3*4))) = uint32(length237)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = uint32(uintptr(uintptr(result236)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length240)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result239)))

switch (payload).StartTime.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (payload).StartTime.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int64((payload).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int32((payload).Nanoseconds)

        
default:
        panic("unreachable")
}
*(*int64)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int32(((payload).Time).Nanoseconds)

        

case MetricDataU64Sum:
        payload := (element).Data.U64Sum()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(5))
slice278 := (payload).DataPoints
length280 := uint32(len(slice278))
result279 := wit_runtime.Allocate(pinner, uintptr(length280 * (16+4*4)), 8)
for index, element := range slice278 {
        base := unsafe.Add(result279, index * (16+4*4))
        slice254 := (element).Attributes
length256 := uint32(len(slice254))
result255 := wit_runtime.Allocate(pinner, uintptr(length256 * (8+4*4)), 8)
for index, element := range slice254 {
        base := unsafe.Add(result255, index * (8+4*4))
        utf8241 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8241)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8241)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8242 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8242)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8242)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result243 int32
if payload {
        result243 = 1
} else {
        result243 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result243)

        

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
slice245 := payload
length247 := uint32(len(slice245))
result246 := wit_runtime.Allocate(pinner, uintptr(length247 * (2*4)), 4)
for index, element := range slice245 {
        base := unsafe.Add(result246, index * (2*4))
        utf8244 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8244)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8244)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length247)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result246)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice249 := payload
length251 := uint32(len(slice249))
result250 := wit_runtime.Allocate(pinner, uintptr(length251 * 1), 1)
for index, element := range slice249 {
        base := unsafe.Add(result250, index * 1)
        var result248 int32
if element {
        result248 = 1
} else {
        result248 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result248)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length251)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result250)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data252 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data252)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data252)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data253 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data253)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data253)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length256)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result255)))

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice275 := (element).Exemplars
length277 := uint32(len(slice275))
result276 := wit_runtime.Allocate(pinner, uintptr(length277 * (32+6*4)), 8)
for index, element := range slice275 {
        base := unsafe.Add(result276, index * (32+6*4))
        slice270 := (element).FilteredAttributes
length272 := uint32(len(slice270))
result271 := wit_runtime.Allocate(pinner, uintptr(length272 * (8+4*4)), 8)
for index, element := range slice270 {
        base := unsafe.Add(result271, index * (8+4*4))
        utf8257 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8257)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8257)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8258 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8258)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8258)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result259 int32
if payload {
        result259 = 1
} else {
        result259 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result259)

        

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
slice261 := payload
length263 := uint32(len(slice261))
result262 := wit_runtime.Allocate(pinner, uintptr(length263 * (2*4)), 4)
for index, element := range slice261 {
        base := unsafe.Add(result262, index * (2*4))
        utf8260 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8260)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8260)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length263)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result262)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice265 := payload
length267 := uint32(len(slice265))
result266 := wit_runtime.Allocate(pinner, uintptr(length267 * 1), 1)
for index, element := range slice265 {
        base := unsafe.Add(result266, index * 1)
        var result264 int32
if element {
        result264 = 1
} else {
        result264 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result264)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length267)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result266)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data268 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data268)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data268)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data269 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data269)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data269)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length272)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result271)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8273 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8273)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8273)))
utf8274 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8274)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8274)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+3*4))) = uint32(length277)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = uint32(uintptr(uintptr(result276)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length280)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result279)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))
var result281 int32
if (payload).IsMonotonic {
        result281 = 1
} else {
        result281 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (41+8*4))) = int8(result281)

        

case MetricDataU64Histogram:
        payload := (element).Data.U64Histogram()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(6))
slice321 := (payload).DataPoints
length323 := uint32(len(slice321))
result322 := wit_runtime.Allocate(pinner, uintptr(length323 * (72+8*4)), 8)
for index, element := range slice321 {
        base := unsafe.Add(result322, index * (72+8*4))
        slice295 := (element).Attributes
length297 := uint32(len(slice295))
result296 := wit_runtime.Allocate(pinner, uintptr(length297 * (8+4*4)), 8)
for index, element := range slice295 {
        base := unsafe.Add(result296, index * (8+4*4))
        utf8282 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8282)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8282)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8283 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8283)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8283)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result284 int32
if payload {
        result284 = 1
} else {
        result284 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result284)

        

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
slice286 := payload
length288 := uint32(len(slice286))
result287 := wit_runtime.Allocate(pinner, uintptr(length288 * (2*4)), 4)
for index, element := range slice286 {
        base := unsafe.Add(result287, index * (2*4))
        utf8285 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8285)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8285)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length288)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result287)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice290 := payload
length292 := uint32(len(slice290))
result291 := wit_runtime.Allocate(pinner, uintptr(length292 * 1), 1)
for index, element := range slice290 {
        base := unsafe.Add(result291, index * 1)
        var result289 int32
if element {
        result289 = 1
} else {
        result289 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result289)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length292)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result291)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data293 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data293)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data293)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data294 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data294)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data294)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length297)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result296)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64((element).Count)
data298 := unsafe.Pointer(unsafe.SliceData((element).Bounds))
pinner.Pin(data298)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len((element).Bounds)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data298)))
data299 := unsafe.Pointer(unsafe.SliceData((element).BucketCounts))
pinner.Pin(data299)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+5*4))) = uint32(uint32(len((element).BucketCounts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+4*4))) = uint32(uintptr(uintptr(data299)))

switch (element).Min.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Min.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Max.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+6*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Max.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+6*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Sum.Tag() {
case MetricNumberF64:
        payload := (element).Sum.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = payload

        

case MetricNumberS64:
        payload := (element).Sum.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = payload

        

case MetricNumberU64:
        payload := (element).Sum.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice318 := (element).Exemplars
length320 := uint32(len(slice318))
result319 := wit_runtime.Allocate(pinner, uintptr(length320 * (32+6*4)), 8)
for index, element := range slice318 {
        base := unsafe.Add(result319, index * (32+6*4))
        slice313 := (element).FilteredAttributes
length315 := uint32(len(slice313))
result314 := wit_runtime.Allocate(pinner, uintptr(length315 * (8+4*4)), 8)
for index, element := range slice313 {
        base := unsafe.Add(result314, index * (8+4*4))
        utf8300 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8300)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8300)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8301 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8301)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8301)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result302 int32
if payload {
        result302 = 1
} else {
        result302 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result302)

        

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
slice304 := payload
length306 := uint32(len(slice304))
result305 := wit_runtime.Allocate(pinner, uintptr(length306 * (2*4)), 4)
for index, element := range slice304 {
        base := unsafe.Add(result305, index * (2*4))
        utf8303 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8303)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8303)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length306)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result305)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice308 := payload
length310 := uint32(len(slice308))
result309 := wit_runtime.Allocate(pinner, uintptr(length310 * 1), 1)
for index, element := range slice308 {
        base := unsafe.Add(result309, index * 1)
        var result307 int32
if element {
        result307 = 1
} else {
        result307 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result307)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length310)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result309)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data311 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data311)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data311)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data312 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data312)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data312)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length315)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result314)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8316 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8316)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8316)))
utf8317 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8317)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8317)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (72+7*4))) = uint32(length320)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (72+6*4))) = uint32(uintptr(uintptr(result319)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length323)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result322)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))

        

case MetricDataU64ExponentialHistogram:
        payload := (element).Data.U64ExponentialHistogram()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(7))
slice363 := (payload).DataPoints
length365 := uint32(len(slice363))
result364 := wit_runtime.Allocate(pinner, uintptr(length365 * (96+10*4)), 8)
for index, element := range slice363 {
        base := unsafe.Add(result364, index * (96+10*4))
        slice337 := (element).Attributes
length339 := uint32(len(slice337))
result338 := wit_runtime.Allocate(pinner, uintptr(length339 * (8+4*4)), 8)
for index, element := range slice337 {
        base := unsafe.Add(result338, index * (8+4*4))
        utf8324 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8324)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8324)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8325 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8325)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8325)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result326 int32
if payload {
        result326 = 1
} else {
        result326 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result326)

        

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
slice328 := payload
length330 := uint32(len(slice328))
result329 := wit_runtime.Allocate(pinner, uintptr(length330 * (2*4)), 4)
for index, element := range slice328 {
        base := unsafe.Add(result329, index * (2*4))
        utf8327 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8327)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8327)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length330)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result329)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice332 := payload
length334 := uint32(len(slice332))
result333 := wit_runtime.Allocate(pinner, uintptr(length334 * 1), 1)
for index, element := range slice332 {
        base := unsafe.Add(result333, index * 1)
        var result331 int32
if element {
        result331 = 1
} else {
        result331 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result331)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length334)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result333)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data335 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data335)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data335)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data336 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data336)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data336)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length339)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result338)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64((element).Count)

switch (element).Min.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Min.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Max.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Max.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Sum.Tag() {
case MetricNumberF64:
        payload := (element).Sum.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Sum.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Sum.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (72+2*4))) = int8(int32((element).Scale))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (80+2*4))) = int64((element).ZeroCount)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (88+2*4))) = ((element).PositiveBucket).Offset
data340 := unsafe.Pointer(unsafe.SliceData(((element).PositiveBucket).Counts))
pinner.Pin(data340)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+4*4))) = uint32(uint32(len(((element).PositiveBucket).Counts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+3*4))) = uint32(uintptr(uintptr(data340)))
*(*int32)(unsafe.Add(unsafe.Pointer(base), (88+5*4))) = ((element).NegativeBucket).Offset
data341 := unsafe.Pointer(unsafe.SliceData(((element).NegativeBucket).Counts))
pinner.Pin(data341)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+7*4))) = uint32(uint32(len(((element).NegativeBucket).Counts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+6*4))) = uint32(uintptr(uintptr(data341)))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (88+8*4))) = (element).ZeroThreshold
slice360 := (element).Exemplars
length362 := uint32(len(slice360))
result361 := wit_runtime.Allocate(pinner, uintptr(length362 * (32+6*4)), 8)
for index, element := range slice360 {
        base := unsafe.Add(result361, index * (32+6*4))
        slice355 := (element).FilteredAttributes
length357 := uint32(len(slice355))
result356 := wit_runtime.Allocate(pinner, uintptr(length357 * (8+4*4)), 8)
for index, element := range slice355 {
        base := unsafe.Add(result356, index * (8+4*4))
        utf8342 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8342)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8342)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8343 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8343)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8343)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result344 int32
if payload {
        result344 = 1
} else {
        result344 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result344)

        

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
slice346 := payload
length348 := uint32(len(slice346))
result347 := wit_runtime.Allocate(pinner, uintptr(length348 * (2*4)), 4)
for index, element := range slice346 {
        base := unsafe.Add(result347, index * (2*4))
        utf8345 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8345)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8345)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length348)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result347)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice350 := payload
length352 := uint32(len(slice350))
result351 := wit_runtime.Allocate(pinner, uintptr(length352 * 1), 1)
for index, element := range slice350 {
        base := unsafe.Add(result351, index * 1)
        var result349 int32
if element {
        result349 = 1
} else {
        result349 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result349)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length352)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result351)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data353 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data353)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data353)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data354 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data354)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data354)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length357)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result356)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8358 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8358)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8358)))
utf8359 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8359)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8359)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (96+9*4))) = uint32(length362)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (96+8*4))) = uint32(uintptr(uintptr(result361)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length365)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result364)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))

        

case MetricDataS64Gauge:
        payload := (element).Data.S64Gauge()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(8))
slice403 := (payload).DataPoints
length405 := uint32(len(slice403))
result404 := wit_runtime.Allocate(pinner, uintptr(length405 * (16+4*4)), 8)
for index, element := range slice403 {
        base := unsafe.Add(result404, index * (16+4*4))
        slice379 := (element).Attributes
length381 := uint32(len(slice379))
result380 := wit_runtime.Allocate(pinner, uintptr(length381 * (8+4*4)), 8)
for index, element := range slice379 {
        base := unsafe.Add(result380, index * (8+4*4))
        utf8366 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8366)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8366)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8367 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8367)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8367)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result368 int32
if payload {
        result368 = 1
} else {
        result368 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result368)

        

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
slice370 := payload
length372 := uint32(len(slice370))
result371 := wit_runtime.Allocate(pinner, uintptr(length372 * (2*4)), 4)
for index, element := range slice370 {
        base := unsafe.Add(result371, index * (2*4))
        utf8369 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8369)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8369)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length372)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result371)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice374 := payload
length376 := uint32(len(slice374))
result375 := wit_runtime.Allocate(pinner, uintptr(length376 * 1), 1)
for index, element := range slice374 {
        base := unsafe.Add(result375, index * 1)
        var result373 int32
if element {
        result373 = 1
} else {
        result373 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result373)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length376)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result375)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data377 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data377)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data377)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data378 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data378)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data378)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length381)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result380)))

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice400 := (element).Exemplars
length402 := uint32(len(slice400))
result401 := wit_runtime.Allocate(pinner, uintptr(length402 * (32+6*4)), 8)
for index, element := range slice400 {
        base := unsafe.Add(result401, index * (32+6*4))
        slice395 := (element).FilteredAttributes
length397 := uint32(len(slice395))
result396 := wit_runtime.Allocate(pinner, uintptr(length397 * (8+4*4)), 8)
for index, element := range slice395 {
        base := unsafe.Add(result396, index * (8+4*4))
        utf8382 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8382)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8382)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8383 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8383)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8383)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result384 int32
if payload {
        result384 = 1
} else {
        result384 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result384)

        

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
slice386 := payload
length388 := uint32(len(slice386))
result387 := wit_runtime.Allocate(pinner, uintptr(length388 * (2*4)), 4)
for index, element := range slice386 {
        base := unsafe.Add(result387, index * (2*4))
        utf8385 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8385)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8385)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length388)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result387)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice390 := payload
length392 := uint32(len(slice390))
result391 := wit_runtime.Allocate(pinner, uintptr(length392 * 1), 1)
for index, element := range slice390 {
        base := unsafe.Add(result391, index * 1)
        var result389 int32
if element {
        result389 = 1
} else {
        result389 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result389)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length392)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result391)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data393 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data393)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data393)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data394 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data394)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data394)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length397)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result396)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8398 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8398)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8398)))
utf8399 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8399)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8399)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+3*4))) = uint32(length402)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = uint32(uintptr(uintptr(result401)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length405)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result404)))

switch (payload).StartTime.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (payload).StartTime.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int64((payload).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int32((payload).Nanoseconds)

        
default:
        panic("unreachable")
}
*(*int64)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int32(((payload).Time).Nanoseconds)

        

case MetricDataS64Sum:
        payload := (element).Data.S64Sum()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(9))
slice443 := (payload).DataPoints
length445 := uint32(len(slice443))
result444 := wit_runtime.Allocate(pinner, uintptr(length445 * (16+4*4)), 8)
for index, element := range slice443 {
        base := unsafe.Add(result444, index * (16+4*4))
        slice419 := (element).Attributes
length421 := uint32(len(slice419))
result420 := wit_runtime.Allocate(pinner, uintptr(length421 * (8+4*4)), 8)
for index, element := range slice419 {
        base := unsafe.Add(result420, index * (8+4*4))
        utf8406 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8406)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8406)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8407 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8407)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8407)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result408 int32
if payload {
        result408 = 1
} else {
        result408 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result408)

        

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
slice410 := payload
length412 := uint32(len(slice410))
result411 := wit_runtime.Allocate(pinner, uintptr(length412 * (2*4)), 4)
for index, element := range slice410 {
        base := unsafe.Add(result411, index * (2*4))
        utf8409 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8409)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8409)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length412)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result411)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice414 := payload
length416 := uint32(len(slice414))
result415 := wit_runtime.Allocate(pinner, uintptr(length416 * 1), 1)
for index, element := range slice414 {
        base := unsafe.Add(result415, index * 1)
        var result413 int32
if element {
        result413 = 1
} else {
        result413 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result413)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length416)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result415)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data417 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data417)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data417)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data418 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data418)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data418)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length421)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result420)))

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice440 := (element).Exemplars
length442 := uint32(len(slice440))
result441 := wit_runtime.Allocate(pinner, uintptr(length442 * (32+6*4)), 8)
for index, element := range slice440 {
        base := unsafe.Add(result441, index * (32+6*4))
        slice435 := (element).FilteredAttributes
length437 := uint32(len(slice435))
result436 := wit_runtime.Allocate(pinner, uintptr(length437 * (8+4*4)), 8)
for index, element := range slice435 {
        base := unsafe.Add(result436, index * (8+4*4))
        utf8422 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8422)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8422)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8423 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8423)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8423)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result424 int32
if payload {
        result424 = 1
} else {
        result424 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result424)

        

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
slice426 := payload
length428 := uint32(len(slice426))
result427 := wit_runtime.Allocate(pinner, uintptr(length428 * (2*4)), 4)
for index, element := range slice426 {
        base := unsafe.Add(result427, index * (2*4))
        utf8425 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8425)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8425)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length428)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result427)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice430 := payload
length432 := uint32(len(slice430))
result431 := wit_runtime.Allocate(pinner, uintptr(length432 * 1), 1)
for index, element := range slice430 {
        base := unsafe.Add(result431, index * 1)
        var result429 int32
if element {
        result429 = 1
} else {
        result429 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result429)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length432)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result431)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data433 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data433)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data433)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data434 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data434)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data434)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length437)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result436)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8438 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8438)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8438)))
utf8439 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8439)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8439)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+3*4))) = uint32(length442)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = uint32(uintptr(uintptr(result441)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length445)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result444)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))
var result446 int32
if (payload).IsMonotonic {
        result446 = 1
} else {
        result446 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (41+8*4))) = int8(result446)

        

case MetricDataS64Histogram:
        payload := (element).Data.S64Histogram()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(10))
slice486 := (payload).DataPoints
length488 := uint32(len(slice486))
result487 := wit_runtime.Allocate(pinner, uintptr(length488 * (72+8*4)), 8)
for index, element := range slice486 {
        base := unsafe.Add(result487, index * (72+8*4))
        slice460 := (element).Attributes
length462 := uint32(len(slice460))
result461 := wit_runtime.Allocate(pinner, uintptr(length462 * (8+4*4)), 8)
for index, element := range slice460 {
        base := unsafe.Add(result461, index * (8+4*4))
        utf8447 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8447)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8447)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8448 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8448)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8448)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result449 int32
if payload {
        result449 = 1
} else {
        result449 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result449)

        

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
slice451 := payload
length453 := uint32(len(slice451))
result452 := wit_runtime.Allocate(pinner, uintptr(length453 * (2*4)), 4)
for index, element := range slice451 {
        base := unsafe.Add(result452, index * (2*4))
        utf8450 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8450)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8450)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length453)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result452)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice455 := payload
length457 := uint32(len(slice455))
result456 := wit_runtime.Allocate(pinner, uintptr(length457 * 1), 1)
for index, element := range slice455 {
        base := unsafe.Add(result456, index * 1)
        var result454 int32
if element {
        result454 = 1
} else {
        result454 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result454)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length457)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result456)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data458 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data458)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data458)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data459 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data459)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data459)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length462)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result461)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64((element).Count)
data463 := unsafe.Pointer(unsafe.SliceData((element).Bounds))
pinner.Pin(data463)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len((element).Bounds)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data463)))
data464 := unsafe.Pointer(unsafe.SliceData((element).BucketCounts))
pinner.Pin(data464)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+5*4))) = uint32(uint32(len((element).BucketCounts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+4*4))) = uint32(uintptr(uintptr(data464)))

switch (element).Min.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Min.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Max.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+6*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Max.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+6*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Sum.Tag() {
case MetricNumberF64:
        payload := (element).Sum.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = payload

        

case MetricNumberS64:
        payload := (element).Sum.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = payload

        

case MetricNumberU64:
        payload := (element).Sum.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+6*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+6*4))) = int64(payload)

        

default:
        panic("unreachable")
}
slice483 := (element).Exemplars
length485 := uint32(len(slice483))
result484 := wit_runtime.Allocate(pinner, uintptr(length485 * (32+6*4)), 8)
for index, element := range slice483 {
        base := unsafe.Add(result484, index * (32+6*4))
        slice478 := (element).FilteredAttributes
length480 := uint32(len(slice478))
result479 := wit_runtime.Allocate(pinner, uintptr(length480 * (8+4*4)), 8)
for index, element := range slice478 {
        base := unsafe.Add(result479, index * (8+4*4))
        utf8465 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8465)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8465)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8466 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8466)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8466)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result467 int32
if payload {
        result467 = 1
} else {
        result467 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result467)

        

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
slice469 := payload
length471 := uint32(len(slice469))
result470 := wit_runtime.Allocate(pinner, uintptr(length471 * (2*4)), 4)
for index, element := range slice469 {
        base := unsafe.Add(result470, index * (2*4))
        utf8468 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8468)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8468)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length471)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result470)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice473 := payload
length475 := uint32(len(slice473))
result474 := wit_runtime.Allocate(pinner, uintptr(length475 * 1), 1)
for index, element := range slice473 {
        base := unsafe.Add(result474, index * 1)
        var result472 int32
if element {
        result472 = 1
} else {
        result472 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result472)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length475)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result474)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data476 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data476)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data476)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data477 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data477)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data477)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length480)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result479)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8481 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8481)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8481)))
utf8482 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8482)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8482)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (72+7*4))) = uint32(length485)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (72+6*4))) = uint32(uintptr(uintptr(result484)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length488)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result487)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))

        

case MetricDataS64ExponentialHistogram:
        payload := (element).Data.S64ExponentialHistogram()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (6*4))) = int8(int32(11))
slice528 := (payload).DataPoints
length530 := uint32(len(slice528))
result529 := wit_runtime.Allocate(pinner, uintptr(length530 * (96+10*4)), 8)
for index, element := range slice528 {
        base := unsafe.Add(result529, index * (96+10*4))
        slice502 := (element).Attributes
length504 := uint32(len(slice502))
result503 := wit_runtime.Allocate(pinner, uintptr(length504 * (8+4*4)), 8)
for index, element := range slice502 {
        base := unsafe.Add(result503, index * (8+4*4))
        utf8489 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8489)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8489)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8490 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8490)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8490)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result491 int32
if payload {
        result491 = 1
} else {
        result491 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result491)

        

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
slice493 := payload
length495 := uint32(len(slice493))
result494 := wit_runtime.Allocate(pinner, uintptr(length495 * (2*4)), 4)
for index, element := range slice493 {
        base := unsafe.Add(result494, index * (2*4))
        utf8492 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8492)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8492)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length495)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result494)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice497 := payload
length499 := uint32(len(slice497))
result498 := wit_runtime.Allocate(pinner, uintptr(length499 * 1), 1)
for index, element := range slice497 {
        base := unsafe.Add(result498, index * 1)
        var result496 int32
if element {
        result496 = 1
} else {
        result496 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result496)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length499)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result498)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data500 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data500)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data500)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data501 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data501)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data501)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length504)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result503)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64((element).Count)

switch (element).Min.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Min.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Max.Tag() {
case wit_types.OptionNone:
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = int8(int32(0))

        
case wit_types.OptionSome:
        payload := (element).Max.Some()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = int8(int32(1))

switch payload.Tag() {
case MetricNumberF64:
        payload := payload.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = payload

        

case MetricNumberS64:
        payload := payload.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = payload

        

case MetricNumberU64:
        payload := payload.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (40+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (48+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}

        
default:
        panic("unreachable")
}

switch (element).Sum.Tag() {
case MetricNumberF64:
        payload := (element).Sum.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Sum.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Sum.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (56+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (64+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (72+2*4))) = int8(int32((element).Scale))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (80+2*4))) = int64((element).ZeroCount)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (88+2*4))) = ((element).PositiveBucket).Offset
data505 := unsafe.Pointer(unsafe.SliceData(((element).PositiveBucket).Counts))
pinner.Pin(data505)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+4*4))) = uint32(uint32(len(((element).PositiveBucket).Counts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+3*4))) = uint32(uintptr(uintptr(data505)))
*(*int32)(unsafe.Add(unsafe.Pointer(base), (88+5*4))) = ((element).NegativeBucket).Offset
data506 := unsafe.Pointer(unsafe.SliceData(((element).NegativeBucket).Counts))
pinner.Pin(data506)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+7*4))) = uint32(uint32(len(((element).NegativeBucket).Counts)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (88+6*4))) = uint32(uintptr(uintptr(data506)))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (88+8*4))) = (element).ZeroThreshold
slice525 := (element).Exemplars
length527 := uint32(len(slice525))
result526 := wit_runtime.Allocate(pinner, uintptr(length527 * (32+6*4)), 8)
for index, element := range slice525 {
        base := unsafe.Add(result526, index * (32+6*4))
        slice520 := (element).FilteredAttributes
length522 := uint32(len(slice520))
result521 := wit_runtime.Allocate(pinner, uintptr(length522 * (8+4*4)), 8)
for index, element := range slice520 {
        base := unsafe.Add(result521, index * (8+4*4))
        utf8507 := unsafe.Pointer(unsafe.StringData((element).Key))
pinner.Pin(utf8507)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len((element).Key)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8507)))

switch (element).Value.Tag() {
case wasi_otel_types.ValueString:
        payload := (element).Value.String()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(0))
utf8508 := unsafe.Pointer(unsafe.StringData(payload))
pinner.Pin(utf8508)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(utf8508)))

        

case wasi_otel_types.ValueBool:
        payload := (element).Value.Bool()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(1))
var result509 int32
if payload {
        result509 = 1
} else {
        result509 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int8(result509)

        

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
slice511 := payload
length513 := uint32(len(slice511))
result512 := wit_runtime.Allocate(pinner, uintptr(length513 * (2*4)), 4)
for index, element := range slice511 {
        base := unsafe.Add(result512, index * (2*4))
        utf8510 := unsafe.Pointer(unsafe.StringData(element))
pinner.Pin(utf8510)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(uint32(len(element)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(utf8510)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length513)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result512)))

        

case wasi_otel_types.ValueBoolArray:
        payload := (element).Value.BoolArray()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(5))
slice515 := payload
length517 := uint32(len(slice515))
result516 := wit_runtime.Allocate(pinner, uintptr(length517 * 1), 1)
for index, element := range slice515 {
        base := unsafe.Add(result516, index * 1)
        var result514 int32
if element {
        result514 = 1
} else {
        result514 = 0
}
*(*int8)(unsafe.Add(unsafe.Pointer(base), 0)) = int8(result514)

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(length517)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(result516)))

        

case wasi_otel_types.ValueF64Array:
        payload := (element).Value.F64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(6))
data518 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data518)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data518)))

        

case wasi_otel_types.ValueS64Array:
        payload := (element).Value.S64Array()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int8(int32(7))
data519 := unsafe.Pointer(unsafe.SliceData(payload))
pinner.Pin(data519)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+3*4))) = uint32(uint32(len(payload)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = uint32(uintptr(uintptr(data519)))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), 4)) = uint32(length522)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), 0)) = uint32(uintptr(uintptr(result521)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (2*4))) = int64(((element).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (8+2*4))) = int32(((element).Time).Nanoseconds)

switch (element).Value.Tag() {
case MetricNumberF64:
        payload := (element).Value.F64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(0))
*(*float64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberS64:
        payload := (element).Value.S64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(1))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = payload

        

case MetricNumberU64:
        payload := (element).Value.U64()
        *(*int8)(unsafe.Add(unsafe.Pointer(base), (16+2*4))) = int8(int32(2))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+2*4))) = int64(payload)

        

default:
        panic("unreachable")
}
utf8523 := unsafe.Pointer(unsafe.StringData((element).SpanId))
pinner.Pin(utf8523)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+3*4))) = uint32(uint32(len((element).SpanId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+2*4))) = uint32(uintptr(uintptr(utf8523)))
utf8524 := unsafe.Pointer(unsafe.StringData((element).TraceId))
pinner.Pin(utf8524)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+5*4))) = uint32(uint32(len((element).TraceId)))
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (32+4*4))) = uint32(uintptr(uintptr(utf8524)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (96+9*4))) = uint32(length527)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (96+8*4))) = uint32(uintptr(uintptr(result526)))

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+7*4))) = uint32(length530)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (8+6*4))) = uint32(uintptr(uintptr(result529)))
*(*int64)(unsafe.Add(unsafe.Pointer(base), (8+8*4))) = int64(((payload).StartTime).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (16+8*4))) = int32(((payload).StartTime).Nanoseconds)
*(*int64)(unsafe.Add(unsafe.Pointer(base), (24+8*4))) = int64(((payload).Time).Seconds)
*(*int32)(unsafe.Add(unsafe.Pointer(base), (32+8*4))) = int32(((payload).Time).Nanoseconds)
*(*int8)(unsafe.Add(unsafe.Pointer(base), (40+8*4))) = int8(int32((payload).Temporality))

        

default:
        panic("unreachable")
}

}

*(*uint32)(unsafe.Add(unsafe.Pointer(base), (11*4))) = uint32(length533)
*(*uint32)(unsafe.Add(unsafe.Pointer(base), (10*4))) = uint32(uintptr(uintptr(result532)))

}

wasm_import_export(uintptr(result9), length10, option, option12, option13, uintptr(result535), length536, returnArea)
var result537 wit_types.Result[wit_types.Unit, string]
switch uint8(*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), 0))) {
case 0:
        
        result537 = wit_types.Ok[wit_types.Unit, string](wit_types.Unit{})
case 1:
        value := unsafe.String((*uint8)(unsafe.Pointer(uintptr(*(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), 4))))), *(*uint32)(unsafe.Add(unsafe.Pointer(returnArea), (2*4))))

        result537 = wit_types.Err[wit_types.Unit, string](value)
default:
        panic("unreachable")
}
result538 := result537;
return result538

}

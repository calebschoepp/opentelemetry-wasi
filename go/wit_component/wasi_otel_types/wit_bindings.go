package wasi_otel_types

import (
        "wit_component/wit_types"
)

// The key part of attribute `key-value` pairs.
type Key = string

const (
// A string value.
ValueString uint8 = 0
// A boolean value.
ValueBool uint8 = 1
// A double precision floating point value.
ValueF64 uint8 = 2
// A signed 64 bit integer value.
ValueS64 uint8 = 3
// A homogeneous array of string values.
ValueStringArray uint8 = 4
// A homogeneous array of boolean values.
ValueBoolArray uint8 = 5
// A homogeneous array of double precision floating point values.
ValueF64Array uint8 = 6
// A homogeneous array of 64 bit integer values.
ValueS64Array uint8 = 7
)

// The value part of attribute `key-value` pairs.
type Value struct {
        tag uint8
        value any
}

func (self Value) Tag() uint8 {
        return self.tag
}

func (self Value) String() string {
        if self.tag != ValueString {
                panic("tag mismatch")
        }
        return self.value.(string)
}
func (self Value) Bool() bool {
        if self.tag != ValueBool {
                panic("tag mismatch")
        }
        return self.value.(bool)
}
func (self Value) F64() float64 {
        if self.tag != ValueF64 {
                panic("tag mismatch")
        }
        return self.value.(float64)
}
func (self Value) S64() int64 {
        if self.tag != ValueS64 {
                panic("tag mismatch")
        }
        return self.value.(int64)
}
func (self Value) StringArray() []string {
        if self.tag != ValueStringArray {
                panic("tag mismatch")
        }
        return self.value.([]string)
}
func (self Value) BoolArray() []bool {
        if self.tag != ValueBoolArray {
                panic("tag mismatch")
        }
        return self.value.([]bool)
}
func (self Value) F64Array() []float64 {
        if self.tag != ValueF64Array {
                panic("tag mismatch")
        }
        return self.value.([]float64)
}
func (self Value) S64Array() []int64 {
        if self.tag != ValueS64Array {
                panic("tag mismatch")
        }
        return self.value.([]int64)
}

func MakeValueString(value string) Value {
        return Value{ValueString, value}
}
func MakeValueBool(value bool) Value {
        return Value{ValueBool, value}
}
func MakeValueF64(value float64) Value {
        return Value{ValueF64, value}
}
func MakeValueS64(value int64) Value {
        return Value{ValueS64, value}
}
func MakeValueStringArray(value []string) Value {
        return Value{ValueStringArray, value}
}
func MakeValueBoolArray(value []bool) Value {
        return Value{ValueBoolArray, value}
}
func MakeValueF64Array(value []float64) Value {
        return Value{ValueF64Array, value}
}
func MakeValueS64Array(value []int64) Value {
        return Value{ValueS64Array, value}
}



// A key-value pair describing an attribute.
type KeyValue struct {
        // The attribute name.
Key string
// The attribute value.
Value Value 
}

// Describes the instrumentation scope that produced telemetry.
type InstrumentationScope struct {
        // Name of the instrumentation scope.
Name string
// The library version.
Version wit_types.Option[string]
// Schema URL used by this library.
// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.9.0/specification/schemas/overview.md#schema-url
SchemaUrl wit_types.Option[string]
// Specifies the instrumentation scope attributes to associate with emitted telemetry.
Attributes []KeyValue 
}

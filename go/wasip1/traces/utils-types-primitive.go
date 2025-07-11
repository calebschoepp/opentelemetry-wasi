//go:build cgo
// +build cgo

package trace

/*
#include <stdint.h>
#include <stdlib.h>
#include "traces.h"
*/
import "C"
import "unsafe"

// Translates a `C.traces_string_t` to a byte slice
func otelStringToGoByteSlice(otelStr C.traces_string_t) []byte {
	if otelStr.ptr == nil || otelStr.len == 0 {
		return nil
	}
	// Direct conversion from C memory to Go slice
	return (*[1 << 28]byte)(unsafe.Pointer(otelStr.ptr))[:otelStr.len:otelStr.len] // TODO: figure out what this does
}

// Translates a `C.traces_string_t` to a string
func otelStringToGoString(otelStr C.traces_string_t) string {
	if otelStr.ptr == nil || otelStr.len == 0 {
		return ""
	}
	return C.GoStringN(otelStr.ptr, C.int(otelStr.len))
}

// Translates from `string` to `C.traces_string_t`
func goStringToOtelString(s string) C.traces_string_t {
	if s == "" {
		return C.traces_string_t{ptr: nil, len: 0}
	}

	cStr := C.CString(s)
	return C.traces_string_t{
		ptr: cStr,
		len: C.size_t(len(s)),
	}
}

func goStringToOtelOptionString(s string) C.traces_option_string_t {
	if s == "" {
		return C.traces_option_string_t{
			is_some: false,
			val:     goStringToOtelString(""),
		}
	}

	return C.traces_option_string_t{
		is_some: true,
		val:     goStringToOtelString(s),
	}
}

func goStringSliceToOtelListString(list []string) C.traces_list_string_t {
	if len(list) == 0 {
		return C.traces_list_string_t{ptr: nil, len: 0}
	}

	cArray := (*C.traces_string_t)(C.malloc(C.size_t(len(list)) * C.size_t(unsafe.Sizeof(C.traces_string_t{}))))
	for i, s := range list {
		*(*C.traces_string_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(C.traces_string_t{}))) = goStringToOtelString(s)
	}

	return C.traces_list_string_t{
		ptr: cArray,
		len: C.size_t(len(list)),
	}
}

func goBoolToCBool(v bool) C.bool {
	if v {
		return C.bool(true)
	}

	return C.bool(false)
}

func goBoolSliceToOtelListBool(list []bool) C.traces_list_bool_t {
	if len(list) == 0 {
		return C.traces_list_bool_t{ptr: nil, len: 0}
	}

	cArray := (*C.bool)(C.malloc(C.size_t(len(list)) * C.size_t(unsafe.Sizeof(C.bool(false)))))
	for i, v := range list {
		*(*C.bool)(unsafe.Pointer(uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(C.bool(false)))) = C.bool(v)
	}

	return C.traces_list_bool_t{
		ptr: cArray,
		len: C.size_t(len(list)),
	}
}

func goSliceF64toOtelListF64(list []float64) C.traces_list_float64_t {
	if len(list) == 0 {
		return C.traces_list_float64_t{ptr: nil, len: 0}
	}

	cArray := (*C.double)(C.malloc(C.size_t(len(list)) * C.size_t(unsafe.Sizeof(C.double(0)))))

	for i, num := range list {
		*(*C.double)(unsafe.Pointer(uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(C.double(0)))) = C.double(num)
	}

	return C.traces_list_float64_t{
		ptr: cArray,
		len: C.size_t(len(list)),
	}
}

func goSliceS64ToOtelListS64(list []int64) C.traces_list_s64_t {
	if len(list) == 0 {
		return C.traces_list_s64_t{ptr: nil, len: 0}
	}

	cArray := (*C.int64_t)(C.malloc(C.size_t(len(list)) * C.size_t(unsafe.Sizeof(C.int64_t(0)))))

	for i, num := range list {
		*(*C.int64_t)(unsafe.Pointer(uintptr(unsafe.Pointer(cArray)) + uintptr(i)*unsafe.Sizeof(C.int64_t(0)))) = C.int64_t(num)
	}

	return C.traces_list_s64_t{
		ptr: cArray,
		len: C.size_t(len(list)),
	}
}

package wasi_clocks_wall_clock

import (
        "wit_component/wit_runtime"
"unsafe"
"runtime"
)


// A time and date in seconds plus nanoseconds.
type Datetime struct {
        Seconds uint64
Nanoseconds uint32 
}

//go:wasmimport wasi:clocks/wall-clock@0.2.0 now
func wasm_import_now(arg0 uintptr) 

func Now() Datetime {
        pinner := &runtime.Pinner{}
defer pinner.Unpin()

        returnArea := uintptr(wit_runtime.Allocate(pinner, 16, 8))
        wasm_import_now(returnArea)
result := Datetime{uint64(*(*int64)(unsafe.Add(unsafe.Pointer(returnArea), 0))), uint32(*(*int32)(unsafe.Add(unsafe.Pointer(returnArea), 8)))};
return result

}

//go:wasmimport wasi:clocks/wall-clock@0.2.0 resolution
func wasm_import_resolution(arg0 uintptr) 

func Resolution() Datetime {
        pinner := &runtime.Pinner{}
defer pinner.Unpin()

        returnArea := uintptr(wit_runtime.Allocate(pinner, 16, 8))
        wasm_import_resolution(returnArea)
result := Datetime{uint64(*(*int64)(unsafe.Add(unsafe.Pointer(returnArea), 0))), uint32(*(*int32)(unsafe.Add(unsafe.Pointer(returnArea), 8)))};
return result

}

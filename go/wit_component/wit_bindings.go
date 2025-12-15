package main

import (
        "runtime"
        "wit_component/wit_runtime"
)

var staticPinner = runtime.Pinner{}
var exportReturnArea = uintptr(wit_runtime.Allocate(&staticPinner, 0, 1))
var syncExportPinner = runtime.Pinner{}



// Unused, but present to make the compiler happy
func main() {}

package main

import (
	"fmt"

	"sdk.com"
)

func init() {
	fmt.Println(sdk.AddExternal(123, 456))
}

func main() {}

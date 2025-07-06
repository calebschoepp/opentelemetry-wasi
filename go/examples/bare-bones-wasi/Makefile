# NOTE: the `tinygo build` command also works with `-target=wasip2`
.phony: build
build:
	@wkg wit fetch -d pkg/sdk/wit
	@wkg wit build -d pkg/sdk/wit
	@cd pkg/app && tinygo build -target=wasip1 -o ../../main.wasm main.go

.phony: run
run:
	@wasmtime main.wasm
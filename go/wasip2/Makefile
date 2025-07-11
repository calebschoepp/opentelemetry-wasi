.phony: build
build:
	@wkg wit fetch
	@wkg wit build

.phony: generate
generate:
	go tool wit-bindgen-go generate --world imports --out internal ./wasi:otel@0.2.0-draft.wasm
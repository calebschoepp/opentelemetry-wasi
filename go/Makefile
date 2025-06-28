.phony: build
build:
	@wkg wit fetch -d ../wit
	@wkg wit build -d ../wit

.phony: generate
generate:
	go tool wit-bindgen-go generate --world imports --out internal ./wasi:otel@0.2.0-draft.wasm
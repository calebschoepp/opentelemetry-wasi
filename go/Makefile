.phony: wkg
wkg:
	@wkg wit fetch -d ../wit
	@wkg wit build -d ../wit

.phony: tool
tool:
	go tool wit-bindgen-go generate --world imports --out internal ./wasi:otel@0.2.0-draft.wasm

.phony: build
build:
	@tinygo build -target=wasip2 -o otel.wasm --wit-package wasi:otel@0.2.0-draft.wasm --wit-world imports main.go



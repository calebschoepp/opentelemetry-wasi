# OpenTelemetry WASI for Go

## Overview
Resources I'm using: 
- https://component-model.bytecodealliance.org/language-support/go.html
- https://opentelemetry.io/docs/languages/go/getting-started/
- https://opentelemetry.io/docs/specs/otel/trace/sdk/

Resources I need to look further-into:
- https://github.com/bytecodealliance/wit-bindgen
- https://github.com/bytecodealliance/go-modules

## Requirements
- [Go](https://go.dev/dl/) version `1.24`
- [Tinygo](https://github.com/tinygo-org/tinygo/releases/tag/v0.38.0) version `0.38`
- [wkg](https://github.com/bytecodealliance/wasm-pkg-tools) version 0.11.0
    - It's simplest to `cargo install wkg`

## Usage

Build a version of [Spin](https://github.com/spinframework/spin) from this [branch](https://github.com/calebschoepp/spin/tree/wasi-otel) and then run the example of your choosing.

```sh
git clone --depth 1 --branch wasi-otel https://github.com/calebschoepp/spin
cargo install --path spin
spin plugin update
spin plugin install otel
```

```sh
cd examples/spin-basic
spin build
spin otel setup
spin otel up
curl localhost:3000
```
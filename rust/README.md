# OpenTelemetry WASI for Rust

Enables using OpenTelemetry within Rust WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

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

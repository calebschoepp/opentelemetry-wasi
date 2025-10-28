# OpenTelemetry WASI for Rust

Enables using OpenTelemetry within Rust WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

## Usage

Build a version of [Spin](https://github.com/spinframework/spin) from this [branch](https://github.com/calebschoepp/spin/tree/wasi-otel) and install the relevant plugins:
```sh
git clone --branch wasi-otel --depth 1 https://github.com/calebschoepp/spin
cd spin
cargo install --path .
spin plugin update
spin plugin install otel
```

Then, run the example of your choosing:
```sh
cd examples/spin-basic
spin build
spin otel setup
spin otel up
# In a different terminal...
curl localhost:3000
```

## Notes about Tracing

## Notes about Metrics
### Observable (Async) Instruments
Async instruments (observable counters, gauges, etc.) collect metric data that must be manually exported to the host. While typical applications use periodic exporters to handle this automatically, Rust WebAssembly applications don't yet support periodic exporters. To address this, this SDK provides a manual reader that will be explicitly called to export the metric data at one or more points during the life of the guest application.

## Notes about Logs
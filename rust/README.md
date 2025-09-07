# OpenTelemetry WASI for Rust

Enables using OpenTelemetry within Rust WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

## Usage

Build a version of [Spin](https://github.com/spinframework/spin) from this [branch](https://github.com/calebschoepp/spin/tree/wasi-otel) and then run the example of your choosing.

```sh
git clone --branch wasi-otel --depth 1 https://github.com/calebschoepp/spin
cd spin
cargo install --path .
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

## Notes about Tracing

## Notes about Metrics

We based the metrics portion of the Rust SDK on [`opentelemetry-prometheus`](https://github.com/open-telemetry/opentelemetry-rust/tree/c811cde1ae21c624870c1b952190e687b16f76b8/opentelemetry-prometheus).

To summarize:
- We initialize a [`ManualReader`](https://github.com/open-telemetry/opentelemetry-rust/blob/c811cde1ae21c624870c1b952190e687b16f76b8/opentelemetry-sdk/src/metrics/manual_reader.rs), and `Arc::clone` the reader into an _**Exporter**_ and a _**Collector**_.
- We pass the _**Exporter**_ into an [`SdkMeterProvider`](https://github.com/open-telemetry/opentelemetry-rust/blob/c811cde1ae21c624870c1b952190e687b16f76b8/opentelemetry-sdk/src/metrics/meter_provider.rs), and when the various instruments associated with the provider are invoked, they will fill the `ManualReader` with metric data.
- When we are ready to pass the metric information from the Wasm guest to the Wasm host, we will use the _**Collector**_ to retrieve a [`ResourceMetrics`](https://github.com/open-telemetry/opentelemetry-rust/blob/c811cde1ae21c624870c1b952190e687b16f76b8/opentelemetry-sdk/src/metrics/data/mod.rs#L13) struct from the `ManualReader` and pass it to the host.

TODO: Add notes about the host implementation

## Notes about Logs

# opentelemetry-wasi

This is a WIP opentelemetry backend based on WASI Observe.

## Usage

Steps to run the sample Spin app are below. Note that you need to be using this custom version of [Spin](https://github.com/fermyon/spin/pull/2485). Also make sure you have an [observability stack](https://github.com/fermyon/spin/tree/main/hack/o11y-stack) running on your system.

```bash
cd examples/exporter/opentelemetry
spin build
OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4318/v1/traces spin up
curl http://localhost:3000
```

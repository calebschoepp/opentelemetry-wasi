# OpenTelemetry WASI for Go

## Usage

### Prerequisites

- [**go**](https://go.dev/dl/) - v1.25+
- [**componentize-go**](https://github.com/bytecodealliance/componentize-go) - v0.2.0
- [**Rust toolchain**](https://rust-lang.org/) - Latest version
- [**Spin**](https://github.com/spinframework/spin) - v3.6.1

### Run an Example Application

```sh
# Setup OTel collector and dashboards
spin plugin update
spin plugin install otel
spin otel setup
spin otel open jaeger # Dashboard for Traces
spin otel open grafana # Dashboard for Metrics and Logs

# Run the application
cd examples/spin-basic
spin otel up -- --build

# Invoke the application (in a different terminal)
curl localhost:3000
```

## Generating the WIT bindings

Whenever WIT files are changed/added to the `../wit` directory, the bindings  in `./wit_component` need to be regenerated.

### Prerequisites

- [**componentize-go**](https://github.com/bytecodealliance/componentize-go) - v0.2.0

### Run

```sh
componentize-go -w imports -d ../wit bindings -o internal --pkg-name github.com/calebschoepp/opentelemetry-wasi/internal --format
```

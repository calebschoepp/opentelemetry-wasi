# OpenTelemetry WASI for TypeScript

Enables using OpenTelemetry within TypeScript WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

## Usage

### Prerequisites

- [**nodejs**](https://nodejs.org/en/download) - Latest version
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

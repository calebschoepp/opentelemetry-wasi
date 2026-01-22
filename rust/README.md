# OpenTelemetry WASI for Rust

Enables using OpenTelemetry within Rust WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

## Usage
### Prerequisites
- [**Rust toolchain**](https://rust-lang.org/) - Latest version
- **Spin** - Installation instructions below:
    ```sh
    # Spin WebAssembly Runtime
    cargo install --git https://github.com/asteurer/spin --rev cc558b2ec0cb2cd0619c6a410325bd165f632f1e spin-cli
    ```

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

## Notes about Metrics
### Observable (Async) Instruments
Async instruments (observable counters, gauges, etc.) collect metric data that must be manually exported to the host. While typical applications use periodic exporters to handle this automatically, Rust WebAssembly applications don't yet support periodic exporters. To address this, this SDK provides a manual reader that will be explicitly called to export the metric data at one or more points during the life of the guest application.

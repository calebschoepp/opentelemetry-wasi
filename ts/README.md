# OpenTelemetry WASI for TypeScript

Enables using OpenTelemetry within TypeScript WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

## Usage
### Prerequisites
- [**nodejs**](https://nodejs.org/en/download) - Latest version
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

# OpenTelemetry WASI for Go
## Usage
### Prerequisites
- [**go**](https://go.dev/dl/) - v1.25+
- [**componentize-go**](https://github.com/asteurer/componentize-go) - Latest version
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

## Generating the WIT bindings
Whenever WIT files are changed/added to the `../wit` directory, the bindings  in `./wit_component` need to be regenerated.

### Prerequisites
- [**componentize-go**](https://github.com/asteurer/componentize-go) - Latest version

### Run
```sh
componentize-go -w imports -d ../wit bindings -o wit_component
```

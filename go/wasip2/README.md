# OpenTelemetry WASI for Go

## Requirements
- [Go](https://go.dev/dl/) version `1.24`
- [Tinygo](https://github.com/tinygo-org/tinygo/releases/tag/v0.38.0) version `0.38`

## Usage

```sh
git clone --depth 1 --branch factor-otel https://github.com/asteurer/otel-spin
cargo install --path otel-spin
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
# OpenTelemetry WASI for TypeScript

Enables using OpenTelemetry within TypeScript WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

## Usage
Build a version of [Spin](https://github.com/spinframework/spin) from this [branch](https://github.com/asteurer/spin/tree/wasi-otel) and install the relevant plugins:
```sh
git clone --branch wasi-otel --depth 1 https://github.com/asteurer/spin
cd spin
cargo install --path .
spin plugin update
spin plugin install otel
```

Build the SDK
```sh
just build-sdk
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

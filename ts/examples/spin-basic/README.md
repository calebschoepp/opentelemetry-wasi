# Usage
Install the TypeScript branch of `opentelemetry-wasi`
```bash
git clone --depth 1 --branch ts https://github.com/asteurer/opentelemetry-wasi.git
```

Install the experimental branch of Spin
```bash
git clone --depth 1 --branch wasi-otel-metrics https://github.com/asteurer/spin.git wasi-otel-spin
```

Build Spin
```bash
cd wasi-otel-spin

# The Spin app take forever to initialize if we don't use a release build
cargo build --release
```

Install required plugins
```bash
spin plugin update

spin plugin install otel
```

Run the application
```bash
cd opentelemetry-wasi/ts/examples/spin-basic

/path/to/wasi-otel-spin build && /path/to/wasi-otel-spin otel up

# In a different terminal...
curl localhost:3000
```
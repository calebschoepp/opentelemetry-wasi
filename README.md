# OpenTelemetry WASI
Libraries to enable using OpenTelemetry within WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

## Integration Tests
### Prerequisites
Before running the integration tests, have the following installed:
- [**just**](https://github.com/casey/just) - Command runner
- [**spin**](https://github.com/asteurer/spin/tree/wasi-otel) - WebAssembly runtime with WASI OTel implemented
- [**Rust Toolchain**](https://rust-lang.org/tools/install/)
- [**OpenSSL development libraries**](https://openssl-library.org/) - `libssl-dev` on Debian/Ubuntu, `openssl-devel` on Fedora
- [**pkg-config**](https://www.freedesktop.org/wiki/Software/pkg-config/) - Available in most package managers

### Running tests
```sh
just test
```

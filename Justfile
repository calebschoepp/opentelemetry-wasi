all: lint format test

lint:
    cargo clippy --manifest-path rust/Cargo.toml --all-targets --all-features -- -D warnings
    cargo clippy --manifest-path integration_tests/Cargo.toml --all-targets --all-features -- -D warnings

format:
    cargo fmt --manifest-path rust/Cargo.toml --all -- --check
    cargo fmt --manifest-path integration_tests/Cargo.toml --all -- --check

test:
    cargo test --manifest-path rust/Cargo.toml
    cargo test --manifest-path integration_tests/Cargo.toml

[working-directory: 'go']
generate-go-bindings:
    wit-bindgen go -w imports --out-dir wit_component ./wit

[working-directory: 'go/examples/spin-basic']
build-go-example:
    @# If the wasip1 reactor file isn't present, retrieve it.
    [ -f "wasi_snapshot_preview1.reactor.wasm" ] || curl -OL https://github.com/bytecodealliance/wasmtime/releases/download/v39.0.1/wasi_snapshot_preview1.reactor.wasm

    @# Build the core wasm module.
    GOOS=wasip1 GOARCH=wasm go build -o main.core.wasm -buildmode=c-shared -ldflags=-checklinkname=0 .

    @# Embed the WIT in the core module.
    wasm-tools component embed -w http-trigger --output main.wit.wasm ../../wit main.core.wasm

    # TODO: we might be able to embed multiple worlds into the main.wit.wasm file.
    # wasm-tools component embed -w imports --output main.wit.wasm ../../wit main.wit.wasm

    @# Create a component from the module, adapting WASI preview1 imports to the component model.
    wasm-tools component new --adapt wasi_snapshot_preview1.reactor.wasm main.wit.wasm --output main.wasm

# TODO: Once wasi:otel is merged into Spin, this needs to be removed.
build-integration-tests-base tag="latest":
    docker build -f integration_tests/Dockerfile.base \
        -t "ghcr.io/calebschoepp/opentelemetry-wasi-integration-tests-base:{{tag}}" .

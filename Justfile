all: lint format test

lint:
    # Linting: Rust SDK
    @cargo clippy --manifest-path rust/Cargo.toml --all-targets --all-features -- -D warnings

    # Linting: Go SDK
    @cd go && GOOS=wasip1 GOARCH=wasm golangci-lint run ./...

    # Linting: Integration tests
    @cargo clippy --manifest-path integration_tests/Cargo.toml --all-targets --all-features -- -D warnings

format:
    # Check formatting: Rust SDK
    @cargo fmt --manifest-path rust/Cargo.toml --all -- --check

    # Check formatting: Go SDK
    @if [ -n "$(gofmt -l ./go)" ]; then \
        echo "The following Go files are not formatted. Run 'gofmt -w ./go':"; \
        gofmt -l ./go; \
        exit 1; \
    fi

    # Check formatting: Integration tests
    @cargo fmt --manifest-path integration_tests/Cargo.toml --all -- --check

test:
    # Test: Rust SDK
    @cargo test --manifest-path rust/Cargo.toml

    # Test: Go SDK
    @cd go \
        && GOOS=wasip1 GOARCH=wasm go test -c -o logs_test.wasm ./logs \
        && wasmtime run --dir=. logs_test.wasm

    # Test: Integration tests
    @cargo test --manifest-path integration_tests/Cargo.toml

# TODO: Once wasi:otel is merged into Spin, this needs to be removed.
build-integration-tests-base tag="latest":
    docker build -f integration_tests/Dockerfile.base \
        -t "ghcr.io/asteurer/opentelemetry-wasi-integration-tests-base:{{tag}}" .

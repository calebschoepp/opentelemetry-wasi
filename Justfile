all: lint format test

lint:
    # Lint: Rust SDK
    @cargo clippy --manifest-path rust/Cargo.toml --all-targets --all-features -- -D warnings

    # Lint: TypeScript SDK
    @npm --prefix ts install
    @npm --prefix ts run lint

    # Lint: Go SDK
    @cd go && GOOS=wasip1 GOARCH=wasm golangci-lint run ./logs/... ./metrics/... ./tracing/... ./types/...

    # Lint: Integration Tests
    @cargo clippy --manifest-path integration_tests/Cargo.toml --all-targets --all-features -- -D warnings

format:
    # Format: Rust SDK
    @cargo fmt --manifest-path rust/Cargo.toml --all -- --check

    # Format: TypeScript SDK
    @npm --prefix ts install
    @npm --prefix ts run format:check

    # Format: Go SDK
    @if [ -n "$(gofmt -l ./go)" ]; then \
        echo "The following Go files are not formatted. Run 'gofmt -w ./go':"; \
        gofmt -l ./go; \
        exit 1; \
    fi

    # Format: Integration tests
    @cargo fmt --manifest-path integration_tests/Cargo.toml --all -- --check

test:
    # Test: Rust SDK
    @cargo test --manifest-path rust/Cargo.toml

    # Test: TypeScript SDK
    @npm --prefix ts install
    @npm --prefix ts run build
    @npm --prefix ts test

    # Test: Go SDK
    @cd go \
        && GOOS=wasip1 GOARCH=wasm go test -ldflags=-checklinkname=0 -c -o logs_test.wasm ./logs \
        && wasmtime run logs_test.wasm

    # Test: Integration tests
    @cargo test --manifest-path integration_tests/Cargo.toml

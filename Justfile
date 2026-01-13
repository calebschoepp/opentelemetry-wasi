all: lint format test

lint:
    # Lint: Rust SDK
    cargo clippy --manifest-path rust/Cargo.toml --all-targets --all-features -- -D warnings

    # Lint: TypeScript SDK
    npm --prefix ts install
    npm --prefix ts run lint

    # Lint: Integration Tests
    cargo clippy --manifest-path integration_tests/Cargo.toml --all-targets --all-features -- -D warnings

format:
    # Format: Rust SDK
    cargo fmt --manifest-path rust/Cargo.toml --all -- --check

    # Format: TypeScript SDK
    npm --prefix ts install
    npm --prefix ts run format:check

    # Format: Integration tests
    cargo fmt --manifest-path integration_tests/Cargo.toml --all -- --check

test:
    # Test: Rust SDK
    cargo test --manifest-path rust/Cargo.toml

    # Test: TypeScript SDK
    npm --prefix ts install
    npm --prefix ts run build
    npm --prefix ts test

    # Test: Integration tests
    cargo test --manifest-path integration_tests/Cargo.toml

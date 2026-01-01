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

# TODO: Once wasi:otel is merged into Spin, this needs to be removed.
build-integration-tests-base tag="latest":
    docker build -f integration_tests/Dockerfile.base \
        -t "ghcr.io/asteurer/opentelemetry-wasi-integration-tests-base:{{tag}}" .

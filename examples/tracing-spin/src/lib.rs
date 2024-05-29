use opentelemetry::trace::TracerProvider as _;
use opentelemetry_wasi::propagation::extract_trace_context;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use spin_sdk::key_value::Store;
use tracing::instrument;
use tracing_subscriber::prelude::*;
use tracing_subscriber::registry;

/// A simple Spin HTTP component.
#[http_component]
fn handle_tracing_spin(_req: Request) -> anyhow::Result<impl IntoResponse> {
    let _otel_guard = init_otel();
    let _trace_context_guard = extract_trace_context();
    let span = tracing::info_span!("ooga booga");
    let _guard = span.enter();
    do_nothing();
    do_kv_work();
    // std::thread::sleep(std::time::Duration::from_secs(10));
    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, Fermyon")
        .build())
}

#[instrument]
fn do_nothing() {
    println!("Doing nothing");
}

#[instrument]
fn do_kv_work() {
    let store = Store::open_default().unwrap();
    store.set("foo", String::from("bar").as_bytes()).unwrap();
    store.set("foo", String::from("bar").as_bytes()).unwrap();
}

fn init_otel() -> ShutdownGuard {
    let exporter = opentelemetry_wasi::exporter::WasiExporter::new();
    let provider_builder =
        opentelemetry_sdk::trace::TracerProvider::builder().with_simple_exporter(exporter);
    let provider = provider_builder.build();
    let _ = opentelemetry::global::set_tracer_provider(provider.clone());
    let tracer = provider.tracer("spin-app-tracing-opentelemetry");
    let otel_tracing_layer = tracing_opentelemetry::layer()
        .with_tracer(tracer)
        .with_threads(false);

    registry().with(otel_tracing_layer).init();

    ShutdownGuard
}

/// An RAII implementation for connection to open telemetry services.
///
/// Shutdown of the open telemetry services will happen on `Drop`.
#[must_use]
pub struct ShutdownGuard;

impl Drop for ShutdownGuard {
    fn drop(&mut self) {
        // Give tracer provider a chance to flush any pending traces.
        opentelemetry::global::shutdown_tracer_provider();
    }
}

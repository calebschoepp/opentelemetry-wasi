use opentelemetry::trace::TracerProvider as _;
use opentelemetry::{global, ContextGuard};
use opentelemetry_wasi::propagation::extract_trace_context;
use opentelemetry_wasi::provider::TracerProvider;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use spin_sdk::key_value::Store;
use tracing::instrument;
use tracing_subscriber::prelude::*;
use tracing_subscriber::registry;

#[http_component]
fn spin_guest_function(_req: Request) -> anyhow::Result<impl IntoResponse> {
    let _otel_guard = init_otel();

    let span = tracing::info_span!("spin_guest_function");
    let _span_guard = span.enter();

    compute_something();
    use_kv_store();

    Ok(Response::builder().status(200).build())
}

#[instrument]
fn compute_something() {
    println!("Computing something...");
    let _x = 5 + 2;
}

#[instrument]
fn use_kv_store() {
    let store = Store::open_default().unwrap();
    compute_something();
    store.get("foo").unwrap();
    store.set("foo", String::from("bar").as_bytes()).unwrap();
}

fn init_otel() -> ShutdownGuard {
    // Get a opentelemetry-wasi tracer provider
    let tracer_provider = TracerProvider::default();

    let tracer = tracer_provider.tracer("spin-app-tracing-opentelemetry");
    let otel_tracing_layer = tracing_opentelemetry::layer()
        .with_tracer(tracer)
        .with_threads(false);
    registry().with(otel_tracing_layer).init();

    // Configure the global singleton tracer provider
    let _ = global::set_tracer_provider(tracer_provider);

    // Propagate the parent trace context into the current context
    let trace_context_guard = extract_trace_context().unwrap();

    ShutdownGuard(trace_context_guard)
}

#[must_use]
pub struct ShutdownGuard(ContextGuard);

impl Drop for ShutdownGuard {
    fn drop(&mut self) {
        // Give tracer provider a chance to flush any pending traces.
        opentelemetry::global::shutdown_tracer_provider();
    }
}

// TODO: This sample app is broken. I think because of some interaction between the processor interface and how rust-tracing uses a tracer.

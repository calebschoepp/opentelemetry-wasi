use opentelemetry::{trace::TracerProvider as _, Context};
use opentelemetry_appender_tracing::layer as logging_layer;
use opentelemetry_sdk::{trace::SdkTracerProvider, Resource};
use opentelemetry_wasi::WasiPropagator;
use spin_sdk::{
    http::{IntoResponse, Request, Response},
    http_component,
    key_value::Store,
};
use tracing::instrument;
use tracing_subscriber::{layer::SubscriberExt, registry, util::SubscriberInitExt};

/// A simple Spin HTTP component.
#[http_component]
fn handle_spin_tracing(_req: Request) -> anyhow::Result<impl IntoResponse> {
    // Set up a tracer using the WASI processor
    let wasi_processor = opentelemetry_wasi::WasiSpanProcessor::new();
    let provider = SdkTracerProvider::builder()
        .with_span_processor(wasi_processor)
        .build();
    let tracer = provider.tracer("tracing-spin");

    // Create a tracing layer with the configured tracer and setup the subscriber
    let tracing_layer = tracing_opentelemetry::layer().with_tracer(tracer);

    // Propagate the context from the Wasm host
    let wasi_propagator = opentelemetry_wasi::TraceContextPropagator::new();
    let _guard = wasi_propagator.extract(&Context::current()).attach();

    main_operation();

    // Set up a LoggerProvider using the WASI log processor.
    let log_resource = Resource::builder().with_service_name("spin-logs").build();
    let processor = opentelemetry_wasi::WasiLogProcessor::new(Some(log_resource.clone()));
    let provider = opentelemetry_sdk::logs::SdkLoggerProvider::builder()
        .with_resource(log_resource)
        .with_log_processor(processor)
        .build();

    // Create a logging layer with the configured logger and setup the subscriber
    let log_layer = logging_layer::OpenTelemetryTracingBridge::new(&provider);

    registry()
        .with(tracing_layer)
        .with(log_layer)
        .try_init()
        .unwrap();

    // Trace something
    main_operation();

    // Log something
    tracing::info!(message = "Hello from Rust Tracing!");

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, Spin!")
        .build())
}

#[instrument(fields(my_attribute = "my-value"))]
fn main_operation() {
    tracing::info!(name: "Main span event", foo = "1");
    child_operation();
}

#[instrument()]
fn child_operation() {
    tracing::info!(name: "Sub span event", bar = 1);
    let store = Store::open_default().unwrap();
    store.set("foo", "bar".as_bytes()).unwrap();
}

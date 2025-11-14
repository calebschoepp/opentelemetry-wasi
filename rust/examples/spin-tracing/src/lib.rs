use opentelemetry::trace::TracerProvider as _;
use opentelemetry::Context;
use opentelemetry_sdk::trace::SdkTracerProvider;
use opentelemetry_wasi::WasiPropagator;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use spin_sdk::key_value::Store;
use tracing::instrument;
use tracing_opentelemetry::OpenTelemetrySpanExt;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::registry;
use tracing_subscriber::util::SubscriberInitExt;

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
    registry().with(tracing_layer).try_init().unwrap();

    main_operation();

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, Fermyon")
        .build())
}

#[instrument(fields(my_attribute = "my-value"))]
fn main_operation() {
    // Propagate the context from the Wasm host
    let wasi_propagator = opentelemetry_wasi::TraceContextPropagator::new();
    if let Err(e) =
        tracing::Span::current().set_parent(wasi_propagator.extract(&Context::current()))
    {
        panic!("{e}");
    };

    tracing::info!(name: "Main span event", foo = "1");
    child_operation();
}

#[instrument()]
fn child_operation() {
    tracing::info!(name: "Sub span event", bar = 1);

    let store = Store::open_default().unwrap();
    store.set("foo", "bar".as_bytes()).unwrap();
}

use opentelemetry::{global, Context};
use opentelemetry_sdk::trace::SdkTracerProvider;
use opentelemetry_wasi::WasiPropagator;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use spin_sdk::key_value::Store;

use opentelemetry::{
    trace::{TraceContextExt, Tracer},
    KeyValue,
};

/// A simple Spin HTTP component.
#[http_component]
fn handle_spin_basic(_req: Request) -> anyhow::Result<impl IntoResponse> {
    // Set up a tracer using the WASI processor
    let wasi_processor = opentelemetry_wasi::WasiSpanProcessor::new();
    let tracer_provider = SdkTracerProvider::builder()
        .with_span_processor(wasi_processor)
        .build();
    global::set_tracer_provider(tracer_provider);
    let tracer = global::tracer("basic-spin");

    // Extract context from the Wasm host
    let wasi_propagator = opentelemetry_wasi::TraceContextPropagator::new();
    let _context_guard = wasi_propagator.extract(&Context::current()).attach();

    // Create some spans and events
    tracer.in_span("main-operation", |cx| {
        let span = cx.span();
        span.set_attribute(KeyValue::new("my-attribute", "my-value"));
        span.add_event(
            "Main span event".to_string(),
            vec![KeyValue::new("foo", "1")],
        );
        tracer.in_span("child-operation", |cx| {
            let span = cx.span();
            span.add_event("Sub span event", vec![KeyValue::new("bar", "1")]);

            let store = Store::open_default().unwrap();
            store.set("foo", "bar".as_bytes()).unwrap();
        });
    });

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, Fermyon")
        .build())
}

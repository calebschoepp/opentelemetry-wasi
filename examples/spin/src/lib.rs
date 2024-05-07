use anyhow::Result;
use opentelemetry::global;
use opentelemetry::trace::Span as _;
use opentelemetry::trace::Tracer as _;
use opentelemetry::KeyValue;
use opentelemetry_wasi::provider::TracerProvider;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use spin_sdk::key_value::Store;

/// A simple Spin HTTP component.
#[http_component]
fn handle_spin(_req: Request) -> Result<impl IntoResponse> {
    configure_o11y();

    // Start a span
    let mut _span = global::tracer("spin").start("guest_span");

    do_kv_work();

    Ok(Response::builder()
        .status(200)
        .body("Hello, Fermyon")
        .build())
    // Span dropped at end of function and so it automatically closes
}

fn configure_o11y() {
    // Get a opentelemetry-wasi tracer provider
    let tracer_provider = TracerProvider::default();

    // Configure the global singleton tracer provider
    global::set_tracer_provider(tracer_provider);
}

fn do_kv_work() {
    let mut span = global::tracer("spin").start("do_kv_work");

    let store = Store::open_default().unwrap();
    store.set("foo", String::from("bar").as_bytes()).unwrap();

    span.set_attribute(KeyValue::new("foo", "bar"));

    span.end();
}

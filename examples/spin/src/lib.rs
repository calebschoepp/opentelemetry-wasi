use anyhow::Result;
use opentelemetry::global;
use opentelemetry::trace::Span as _;
use opentelemetry::trace::Tracer as _;
use opentelemetry::KeyValue;
// use opentelemetry_wasi::processor::WasiProcessor;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use spin_sdk::key_value::Store;

/// A simple Spin HTTP component.
#[http_component]
fn handle_spin(_req: Request) -> Result<impl IntoResponse> {
    configure_o11y();

    // Start a span
    let mut span = global::tracer("spin").start("guest_span");
    do_kv_work();

    span.end();
    Ok(Response::builder()
        .status(200)
        .body("Hello, Fermyon")
        .build())
}

fn configure_o11y() {
    let exporter = opentelemetry_wasi::exporter::WasiExporter::new();
    let provider_builder =
        opentelemetry_sdk::trace::TracerProvider::builder().with_simple_exporter(exporter);
    let provider = provider_builder.build();

    let _ = global::set_tracer_provider(provider);
}

fn do_kv_work() {
    let mut span = global::tracer("spin").start("do_kv_work");

    let store = Store::open_default().unwrap();
    store.set("foo", String::from("bar").as_bytes()).unwrap();

    span.set_attribute(KeyValue::new("foo", "bar"));

    span.end();
}

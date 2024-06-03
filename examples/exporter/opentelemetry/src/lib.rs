use anyhow::Result;
use opentelemetry::trace::{TraceContextExt as _, Tracer as _};
use opentelemetry::ContextGuard;
use opentelemetry::{global, Context};
use opentelemetry_wasi::propagation::extract_trace_context;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use spin_sdk::key_value::Store;

#[http_component]
fn spin_guest_function(_req: Request) -> Result<impl IntoResponse> {
    let _otel_guard = init_otel();

    let _guard =
        Context::current_with_span(global::tracer("spin").start("spin_guest_function")).attach();

    compute_something();
    use_kv_store();

    Ok(Response::builder().status(200).build())
}

fn compute_something() {
    let _guard =
        Context::current_with_span(global::tracer("spin").start("compute_something")).attach();
    println!("Computing something...");
    let _x = 5 + 2;
}

fn use_kv_store() {
    let _guard = Context::current_with_span(global::tracer("spin").start("use_kv_store")).attach();
    let store = Store::open_default().unwrap();
    compute_something();
    store.get("foo").unwrap();
    store.set("foo", String::from("bar").as_bytes()).unwrap();
}

fn init_otel() -> ContextGuard {
    let exporter = opentelemetry_wasi::exporter::WasiExporter::new();
    let provider_builder =
        opentelemetry_sdk::trace::TracerProvider::builder().with_simple_exporter(exporter);
    let provider = provider_builder.build();

    let _ = global::set_tracer_provider(provider);

    extract_trace_context().unwrap()
}

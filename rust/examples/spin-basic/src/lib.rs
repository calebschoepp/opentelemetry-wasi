use std::time::Duration;

use opentelemetry::{global, Context};
use opentelemetry_otlp::{Protocol, WithExportConfig};
use opentelemetry_sdk::metrics::SdkMeterProvider;
use opentelemetry_sdk::trace::TracerProvider;
use opentelemetry_sdk::Resource;
use opentelemetry_wasi::WasiPropagator;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use spin_sdk::key_value::Store;

use opentelemetry::{
    trace::{TraceContextExt, Tracer},
    KeyValue,
};

fn init_meter_provider() -> opentelemetry_sdk::metrics::SdkMeterProvider {
    let exporter = opentelemetry_otlp::MetricExporter::builder()
        .with_http()
        .with_endpoint("http://localhost:4318/v1/metrics")
        .with_protocol(Protocol::HttpBinary)
        .with_timeout(Duration::from_secs(3))
        .build()
        .unwrap();

    let reader = opentelemetry_sdk::metrics::PeriodicReader::builder(
        exporter,
        opentelemetry_sdk::runtime::Tokio,
    )
    .with_interval(std::time::Duration::from_secs(3))
    .with_timeout(Duration::from_secs(10))
    .build();

    let provider = opentelemetry_sdk::metrics::SdkMeterProvider::builder()
        .with_reader(reader)
        .with_resource(Resource::new(vec![KeyValue::new(
            "service.name",
            "example",
        )]))
        .build();

    global::set_meter_provider(provider.clone());
    provider
}

/// A simple Spin HTTP component.
#[http_component]
fn handle_spin_basic(_req: Request) -> anyhow::Result<impl IntoResponse> {
    // Try to do metric things
    let meter_provider = init_meter_provider();
    let meter = global::meter("my_service");
    let counter = meter.u64_counter("my_counter").build();
    counter.add(1, &[KeyValue::new("http.client_ip", "83.164.160.102")]);

    // Set up a tracer using the WASI processor
    let wasi_processor = opentelemetry_wasi::WasiProcessor::new();
    let tracer_provider = TracerProvider::builder()
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

    meter_provider.shutdown()?;
    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, Fermyon")
        .build())
}

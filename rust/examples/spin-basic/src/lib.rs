use opentelemetry::{
    global,
    logs::{LogRecord, Logger, LoggerProvider, Severity},
    trace::{TraceContextExt, Tracer},
    Context, KeyValue,
};
use opentelemetry_sdk::{
    logs::SdkLoggerProvider, metrics::SdkMeterProvider, trace::SdkTracerProvider, Resource,
};
use opentelemetry_wasi::{
    TraceContextPropagator, WasiLogProcessor, WasiMetricExporter, WasiPropagator, WasiSpanProcessor,
};
use spin_sdk::{
    http::{IntoResponse, Request, Response},
    http_component,
    key_value::Store,
};

#[http_component]
fn handle_spin_basic(_req: Request) -> anyhow::Result<impl IntoResponse> {
    // ---------------
    // --- Tracing ---
    // ---------------

    // Set up a tracer using the WASI span processor.
    let span_processor = WasiSpanProcessor::new();
    let tracer_provider = SdkTracerProvider::builder()
        .with_span_processor(span_processor)
        .build();
    global::set_tracer_provider(tracer_provider);
    let tracer = global::tracer("basic-spin");

    // Extract context from the Wasm host
    let wasi_propagator = TraceContextPropagator::new();
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

    // ---------------
    // --- Metrics ---
    // ---------------

    // Set up a meter provider using the WASI metric exporter.
    // By default `WasiMetricExporter` will export metrics to the host once it is dropped.
    let metric_exporter = WasiMetricExporter::default();
    let metric_provider = SdkMeterProvider::builder()
        .with_reader(metric_exporter.clone())
        .with_resource(
            Resource::builder()
                .with_service_name("spin-metrics")
                .build(),
        )
        .build();
    global::set_meter_provider(metric_provider.clone());
    let meter = global::meter("spin-meter");

    let attrs = &[
        KeyValue::new("spinkey1", "spinvalue1"),
        KeyValue::new("spinkey2", "spinvalue2"),
    ];

    // Create some instruments and measure things.
    let counter = meter.u64_counter("spin-counter").build();
    counter.add(10, attrs);

    let up_down_counter = meter.i64_up_down_counter("spin-up-down-counter").build();
    up_down_counter.add(-1, attrs);

    let histogram = meter.u64_histogram("spin-histogram").build();
    histogram.record(9, attrs);
    histogram.record(15, attrs);

    let gauge = meter.f64_gauge("spin-gauge").build();
    gauge.record(123.456, attrs);

    // ------------
    // --- Logs ---
    // ------------

    // Set up a LoggerProvider using the WASI log processor.
    let log_resource = Resource::builder().with_service_name("spin-logs").build();
    let log_processor = WasiLogProcessor::new(Some(log_resource.clone()));
    let log_provider = SdkLoggerProvider::builder()
        .with_resource(log_resource)
        .with_log_processor(log_processor)
        .build();
    let logger = log_provider.logger("spin-logger");

    // Create and emit a log.
    let mut record = logger.create_log_record();
    record.set_body("Metrics and traces and logs, oh my!".into());
    record.set_severity_number(Severity::Info);
    record.set_severity_text("info");
    logger.emit(record);

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, world!")
        .build())
}

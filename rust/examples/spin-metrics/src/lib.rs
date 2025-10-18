use opentelemetry::{global, KeyValue};
use opentelemetry_sdk::{metrics::SdkMeterProvider, Resource};
use opentelemetry_wasi::WasiMetricExporter;
use spin_sdk::{
    http::{IntoResponse, Request, Response},
    http_component,
};

#[http_component]
fn handle_spin_metrics(_req: Request) -> anyhow::Result<impl IntoResponse> {
    let reader = WasiMetricExporter::default();
    let provider = SdkMeterProvider::builder()
        .with_reader(reader.clone())
        .with_resource(
            Resource::builder()
                .with_service_name("spin-metrics")
                .build(),
        )
        .build();
    global::set_meter_provider(provider.clone());

    // WARNING: Async instruments (i.e. Observable counters, gauges, etc.) are
    // not yet supported, and will generate a runtime panic.
    let meter = global::meter("spin_meter");
    let counter = meter.u64_counter("spin_counter").build();
    let up_down_counter = meter.i64_up_down_counter("spin_up_down_counter").build();
    let histogram = meter.u64_histogram("spin_histogram").build();
    let gauge = meter.u64_gauge("spin_gauge").build();

    let attrs = &[
        KeyValue::new("spinkey1", "spinvalue1"),
        KeyValue::new("spinkey2", "spinvalue2"),
    ];

    counter.add(10, attrs);
    up_down_counter.add(-1, attrs);
    histogram.record(9, attrs);
    gauge.record(8, attrs);

    reader.export()?;

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, World!")
        .build())
}

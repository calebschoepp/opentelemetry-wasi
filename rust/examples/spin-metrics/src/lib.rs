use opentelemetry::{global, KeyValue};
use opentelemetry_sdk::{
    metrics::{ManualReader, SdkMeterProvider},
    Resource,
};
use opentelemetry_wasi::{WasiMetricCollector, WasiMetricExporter};
use spin_sdk::{
    http::{IntoResponse, Request, Response},
    http_component,
};
use std::sync::Arc;

#[http_component]
fn handle_spin_metrics(_req: Request) -> anyhow::Result<impl IntoResponse> {
    let reader = Arc::new(ManualReader::builder().build());
    let exporter = WasiMetricExporter::new(Arc::clone(&reader));
    let collector = WasiMetricCollector::new(Arc::clone(&reader));
    let provider = SdkMeterProvider::builder()
        .with_reader(exporter)
        .with_resource(
            Resource::builder()
                .with_service_name("spin-metrics")
                .build(),
        )
        .build();
    global::set_meter_provider(provider);

    let counter = global::meter("spin_meter")
        .u64_counter("spin_counter")
        .build();
    let up_down_counter = global::meter("spin_meter")
        .i64_up_down_counter("spin_up_down_counter")
        .build();
    let histogram = global::meter("spin_meter")
        .u64_histogram("spin_histogram")
        .build();
    let gauge = global::meter("spin_meter").u64_gauge("spin_gauge").build();

    counter.add(
        10,
        &[
            KeyValue::new("counterkey1", "countervalue1"),
            KeyValue::new("counterkey2", "countervalue2"),
        ],
    );

    up_down_counter.add(
        -1,
        &[
            KeyValue::new("updowncounterkey1", "updowncountervalue1"),
            KeyValue::new("updowncounterkey2", "updowncountervalue2"),
        ],
    );

    histogram.record(
        9,
        &[
            KeyValue::new("histogramkey1", "histogramvalue1"),
            KeyValue::new("histogramkey2", "histogramvalue2"),
        ],
    );

    gauge.record(
        8,
        &[
            KeyValue::new("gaugekey1", "gaugevalue1"),
            KeyValue::new("gaugekey2", "gaugevalue2"),
        ],
    );

    collector.collect()?;

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, World!")
        .build())
}

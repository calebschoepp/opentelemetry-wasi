use once_cell::sync::Lazy;
use opentelemetry::{
    global,
    metrics::{Counter, Gauge, Histogram, UpDownCounter},
    KeyValue,
};
use opentelemetry_sdk::{
    metrics::{ManualReader, SdkMeterProvider},
    Resource,
};
use opentelemetry_wasi::WasiMetricCollector;
use spin_sdk::{
    http::{IntoResponse, Request, Response},
    http_component,
};
use std::sync::Arc;

static READER: Lazy<Arc<ManualReader>> = Lazy::new(|| Arc::new(ManualReader::builder().build()));

static _PROVIDER: Lazy<()> = Lazy::new(|| {
    Lazy::force(&READER);
    let exporter = opentelemetry_wasi::WasiMetricExporter::new(Arc::clone(&READER));
    let provider = SdkMeterProvider::builder()
        .with_reader(exporter)
        .with_resource(
            Resource::builder()
                .with_service_name("spin-metrics")
                .build(),
        )
        .build();

    global::set_meter_provider(provider);
});

static COLLECTOR: Lazy<WasiMetricCollector> = Lazy::new(|| {
    Lazy::force(&READER);
    opentelemetry_wasi::WasiMetricCollector::new(Arc::clone(&READER))
});

static COUNTER: Lazy<Counter<u64>> = Lazy::new(|| {
    Lazy::force(&_PROVIDER);
    global::meter("spin_meter")
        .u64_counter("spin_counter")
        .build()
});

static UPDOWNCOUNTER: Lazy<UpDownCounter<i64>> = Lazy::new(|| {
    Lazy::force(&_PROVIDER);
    global::meter("spin_meter")
        .i64_up_down_counter("spin_up_down_counter")
        .build()
});

static HISTOGRAM: Lazy<Histogram<u64>> = Lazy::new(|| {
    Lazy::force(&_PROVIDER);
    global::meter("spin_meter")
        .u64_histogram("spin_histogram")
        .build()
});

static GAUGE: Lazy<Gauge<u64>> = Lazy::new(|| {
    Lazy::force(&_PROVIDER);
    global::meter("spin_meter").u64_gauge("spin_gauge").build()
});

#[http_component]
fn handle_spin_metrics(_req: Request) -> anyhow::Result<impl IntoResponse> {
    COUNTER.add(
        10,
        &[
            KeyValue::new("counterkey1", "countervalue1"),
            KeyValue::new("counterkey2", "countervalue2"),
        ],
    );

    UPDOWNCOUNTER.add(
        -1,
        &[
            KeyValue::new("updowncounterkey1", "updowncountervalue1"),
            KeyValue::new("updowncounterkey2", "updowncountervalue2"),
        ],
    );

    HISTOGRAM.record(
        9,
        &[
            KeyValue::new("histogramkey1", "histogramvalue1"),
            KeyValue::new("histogramkey2", "histogramvalue2"),
        ],
    );

    GAUGE.record(
        8,
        &[
            KeyValue::new("gaugekey1", "gaugevalue1"),
            KeyValue::new("gaugekey2", "gaugevalue2"),
        ],
    );

    COLLECTOR.collect()?;

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello, World!")
        .build())
}

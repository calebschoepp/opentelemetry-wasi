use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;

use opentelemetry::{global, KeyValue};
use opentelemetry_sdk::metrics::SdkMeterProvider;
use opentelemetry_sdk::Resource;

/// A simple Spin HTTP component.
#[http_component]
fn handle_spin_metrics(req: Request) -> anyhow::Result<impl IntoResponse> {
    // Initialize the MeterProvider with the stdout Exporter.
    let meter_provider = init_meter_provider();
    global::set_meter_provider(meter_provider);

    // Create a meter from the above MeterProvider.
    let meter = global::meter("mylibraryname");

    // Create a Counter Instrument.
    let counter = meter.u64_counter("my_counter").build();

    // Record measurements using the Counter instrument.
    counter.add(
        10,
        &[
            KeyValue::new("mykey1", "myvalue1"),
            KeyValue::new("mykey2", "myvalue2"),
        ],
    );

    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello World!")
        .build())
}

fn init_meter_provider() -> opentelemetry_sdk::metrics::SdkMeterProvider {
    let exporter = opentelemetry_wasi::WasiExporter{};
    // let exporter = opentelemetry_stdout::MetricExporterBuilder::default()
    //     // Build exporter using Delta Temporality (Defaults to Temporality::Cumulative)
    //     // .with_temporality(opentelemetry_sdk::metrics::Temporality::Delta)
    //     .build();
    let provider = SdkMeterProvider::builder()
        .with_periodic_exporter(exporter)
        .with_resource(
            Resource::builder()
                .with_service_name("metrics-basic-example")
                .build(),
        )
        .build();
    global::set_meter_provider(provider.clone());
    provider
}
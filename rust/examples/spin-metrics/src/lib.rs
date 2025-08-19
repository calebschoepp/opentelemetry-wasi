use opentelemetry::{global, KeyValue};
use opentelemetry_sdk::metrics::{
    periodic_reader_with_async_runtime::PeriodicReader, SdkMeterProvider,
};
use opentelemetry_sdk::{runtime, Resource};
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;

#[http_component]
fn handle_spin_metrics(_req: Request) -> anyhow::Result<impl IntoResponse> {
    let meter_provider = init_meter_provider();
    global::set_meter_provider(meter_provider);

    let meter = global::meter("spin_meter");
    let counter = meter.u64_counter("spin_counter").build();

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
        .body("Hello, World!")
        .build())
}

fn init_meter_provider() -> opentelemetry_sdk::metrics::SdkMeterProvider {
    let exporter = opentelemetry_wasi::WasiMetricExporter::new();
    let reader = PeriodicReader::builder(exporter, runtime::Tokio).build();
    let provider = SdkMeterProvider::builder()
        .with_reader(reader)
        .with_resource(
            Resource::builder()
                .with_service_name("spin-metrics")
                .build(),
        )
        .build();
    global::set_meter_provider(provider.clone());
    provider
}

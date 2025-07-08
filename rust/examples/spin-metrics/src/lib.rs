use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use opentelemetry::metrics::{Counter, Histogram, Meter, MeterProvider};
use opentelemetry::KeyValue;
use opentelemetry_sdk::metrics::{SdkMeterProvider, PeriodicReader, MeterProviderBuilder};
use opentelemetry_sdk::{resource, Resource};

/// A simple Spin HTTP component.
#[http_component]
fn handle_spin_metrics(req: Request) -> anyhow::Result<impl IntoResponse> {
    println!("Handling request to {:?}", req.header("spin-full-url"));
    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello World!")
        .build())
}

fn test() {
    let meter_provider = SdkMeterProvider::builder()
        .with_resource(Resource::new(vec![KeyValue::new("service.name", "my-service")]))
        .build();

    let meter = meter_provider.meter("my-meter");

    let counter= meter
        .u64_counter("http_requests_total")
        .with_description("Total number of HTTP requests")
        .with_unit("requests")
        .build();

    counter.add(15, &[KeyValue::new("status", "200")]);
}
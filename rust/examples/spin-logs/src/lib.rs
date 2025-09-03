use opentelemetry_appender_tracing::layer;
use opentelemetry_sdk::Resource;
use spin_sdk::http::{IntoResponse, Request, Response};
use spin_sdk::http_component;
use tracing_subscriber::prelude::*;

#[http_component]
fn handle_spin_logs(_req: Request) -> anyhow::Result<impl IntoResponse> {
    let processor = opentelemetry_wasi::WasiLogProcessor::new();
    let provider = opentelemetry_sdk::logs::SdkLoggerProvider::builder()
        .with_resource(Resource::builder().with_service_name("spin-logs").build())
        .with_log_processor(processor)
        .build();

    let layer = layer::OpenTelemetryTracingBridge::new(&provider);
    tracing_subscriber::registry().with(layer).init();

    tracing::info!(message = "Hello, world!");

    let _ = provider.shutdown();
    Ok(Response::builder()
        .status(200)
        .header("content-type", "text/plain")
        .body("Hello World!")
        .build())
}

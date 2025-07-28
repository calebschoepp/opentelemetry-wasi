use std::sync::atomic::{Ordering, AtomicBool};
use async_trait::async_trait;

use opentelemetry_sdk::{
    error::{OTelSdkResult, OTelSdkError},
    metrics::{data::ResourceMetrics, exporter::PushMetricExporter, Temporality}
};

use crate::wit::wasi::otel::metrics::export;

pub struct WasiExporter {
    // TODO: I'm just copying what was done in the `tracing/processor.rs` file, so this
    // and the non-export methods implemented on the struct may not be correct...
    pub is_shutdown: AtomicBool,
}

#[async_trait]
impl PushMetricExporter for WasiExporter {
    async fn export(&self, metrics: &mut ResourceMetrics) -> OTelSdkResult {
        let converted = metrics.into();

        export(&converted).map_err(|e| OTelSdkError::InternalFailure(e.to_string()))
    }

    async fn force_flush(&self) -> OTelSdkResult {
        if self.is_shutdown.load(Ordering::Relaxed) {
            return OTelSdkResult::Err(opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown)
        }
        Ok(())
    }

    fn shutdown(&self) -> OTelSdkResult {
        let mut result: Result<(), opentelemetry_sdk::error::OTelSdkError> = Ok(());

        async {result = self.force_flush().await}; // TODO: this might be a no-op...

        if self.is_shutdown.swap(true, Ordering::Relaxed) {
            return OTelSdkResult::Err(opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown)
        }

        result
    }

    fn temporality(&self) -> Temporality{
        Temporality::Cumulative
    }
}
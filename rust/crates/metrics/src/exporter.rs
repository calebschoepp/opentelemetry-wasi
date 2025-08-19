use async_trait::async_trait;
use std::sync::atomic::{AtomicBool, Ordering};

use opentelemetry_sdk::{
    error::{OTelSdkError, OTelSdkResult},
    metrics::{data::ResourceMetrics, exporter::PushMetricExporter, Temporality},
};

use crate::wit::wasi::otel::metrics::export;

pub struct WasiMetricExporter {
    // TODO: I'm just copying what was done in the `tracing/processor.rs` file, so this
    // and the non-export methods implemented on the struct may not be correct...
    pub is_shutdown: AtomicBool,
}

impl WasiMetricExporter {
    pub fn new() -> Self {
        Self {
            is_shutdown: AtomicBool::new(false),
        }
    }
}

#[async_trait]
impl PushMetricExporter for WasiMetricExporter {
    async fn export(&self, metrics: &mut ResourceMetrics) -> OTelSdkResult {
        let converted = metrics.into();

        export(&converted).map_err(|e| OTelSdkError::InternalFailure(e.to_string()))
    }

    async fn force_flush(&self) -> OTelSdkResult {
        if self.is_shutdown.load(Ordering::Relaxed) {
            return OTelSdkResult::Err(opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown);
        }
        Ok(())
    }

    fn shutdown(&self) -> OTelSdkResult {
        if self.is_shutdown.swap(true, Ordering::Relaxed) {
            return OTelSdkResult::Err(opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown);
        }

        Ok(())
    }

    fn temporality(&self) -> Temporality {
        Temporality::Cumulative
    }
}

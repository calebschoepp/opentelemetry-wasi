use std::time::Duration;

use opentelemetry_sdk::{
    error::OTelSdkResult,
    metrics::{data::ResourceMetrics, exporter::PushMetricExporter, Temporality}
};
use crate::wit::wasi::otel::metrics::{export, force_flush, shutdown, temporality, OtelSdkError, TemporalityT};
use async_trait::async_trait;
pub struct WasiExporter {}

#[async_trait]
impl PushMetricExporter for WasiExporter {
    async fn export(&self, metrics: &mut ResourceMetrics) -> OTelSdkResult {
        let converted = metrics.into();
        export(&converted).map_err(|e| match e {
            OtelSdkError::AlreadyShutdown => opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown,
            OtelSdkError::InternalFailure(e) => opentelemetry_sdk::error::OTelSdkError::InternalFailure(e),
            OtelSdkError::Timeout(d) => opentelemetry_sdk::error::OTelSdkError::Timeout(Duration::from_nanos(d)),
        })
    }

    async fn force_flush(&self) -> OTelSdkResult{
        force_flush().map_err(|e| match e {
            OtelSdkError::AlreadyShutdown => opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown,
            OtelSdkError::InternalFailure(e) => opentelemetry_sdk::error::OTelSdkError::InternalFailure(e),
            OtelSdkError::Timeout(d) => opentelemetry_sdk::error::OTelSdkError::Timeout(Duration::from_nanos(d)),
        })
    }

    fn shutdown(&self) -> OTelSdkResult{
        shutdown().map_err(|e| match e {
            OtelSdkError::AlreadyShutdown => opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown,
            OtelSdkError::InternalFailure(e) => opentelemetry_sdk::error::OTelSdkError::InternalFailure(e),
            OtelSdkError::Timeout(d) => opentelemetry_sdk::error::OTelSdkError::Timeout(Duration::from_nanos(d)),
        })
    }

    fn temporality(&self) -> Temporality{
        match temporality(){
            TemporalityT::Cumulative => Temporality::Cumulative,
            TemporalityT::Delta => Temporality::Delta,
            _ => Temporality::LowMemory,
        }
    }
}
use opentelemetry_sdk::metrics::{data::ResourceMetrics, exporter::PushMetricExporter, MetricResult, Temporality};
use crate::wit::wasi::otel::metrics::{export, force_flush, shutdown, temporality};
use async_trait::async_trait;
pub struct WasiExporter {}

#[async_trait]
impl PushMetricExporter for WasiExporter {
    async fn export(&self, metrics: &mut ResourceMetrics) -> MetricResult<()>{
        Ok(())
    }

    async fn force_flush(&self) -> MetricResult<()>{
        Ok(())
    }

    fn shutdown(&self) -> MetricResult<()>{
        Ok(())
    }

    fn temporality(&self) -> Temporality{
        Temporality::LowMemory
    }
}
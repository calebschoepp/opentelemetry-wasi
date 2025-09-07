use crate::wit::wasi::otel::metrics::collect;
use opentelemetry_sdk::{
    error::OTelSdkResult,
    metrics::{
        data::ResourceMetrics, reader::MetricReader, InstrumentKind, ManualReader, Temporality,
    },
};
use std::sync::Arc;

#[derive(Debug)]
pub struct WasiMetricExporter {
    reader: Arc<ManualReader>,
}

impl WasiMetricExporter {
    // This reader will be shared between the collector and the exporter
    pub fn new(reader: Arc<ManualReader>) -> Self {
        Self { reader }
    }
}

impl MetricReader for WasiMetricExporter {
    fn register_pipeline(&self, pipeline: std::sync::Weak<opentelemetry_sdk::metrics::Pipeline>) {
        self.reader.register_pipeline(pipeline)
    }

    fn collect(&self, rm: &mut ResourceMetrics) -> OTelSdkResult {
        self.reader.collect(rm)
    }

    fn force_flush(&self) -> OTelSdkResult {
        self.reader.force_flush()
    }

    fn shutdown(&self) -> OTelSdkResult {
        self.reader.shutdown()
    }

    fn shutdown_with_timeout(&self, timeout: std::time::Duration) -> OTelSdkResult {
        self.reader.shutdown_with_timeout(timeout)
    }

    fn temporality(&self, _kind: InstrumentKind) -> Temporality {
        // TODO: investigate whether we need to use the other Temporalities
        Temporality::Cumulative
    }
}

pub struct WasiMetricCollector {
    reader: Arc<ManualReader>,
}

impl WasiMetricCollector {
    // This reader will be shared between the collector and the exporter
    pub fn new(reader: Arc<ManualReader>) -> Self {
        Self { reader }
    }

    pub fn collect(&self) -> OTelSdkResult {
        // Scrape the metrics from the reader
        let mut metrics = ResourceMetrics::default();
        self.reader.collect(&mut metrics)?;
        collect(&metrics.into()).map_err(|e| e.into())
    }
}

use crate::wit::wasi;
use opentelemetry_sdk::{
    error::OTelSdkResult,
    metrics::{
        data::ResourceMetrics, reader::MetricReader, InstrumentKind, ManualReader, Temporality,
    },
};
use std::sync::Arc;

/// A metric exporter that sends OpenTelemetry metrics to a WASI host.
///
/// This exporter wraps a `ManualReader` to use the reader's existing OpenTelemetry SDK internals
/// while providing a WASI export mechanism. The embedded reader handles metric
/// collection and aggregation, and the wrapper manages exports to the host.
///
/// # Example
/// ```ignore
/// let exporter = WasiMetricExporter::default();
/// let provider = SdkMeterProvider::builder().with_reader(reader.clone());
/// // Measure something...
/// exporter.export()?; // Manually trigger export to WASI host
/// ```
#[derive(Debug, Clone)]
pub struct WasiMetricExporter {
    reader: Arc<ManualReader>,
}

impl Default for WasiMetricExporter {
    fn default() -> Self {
        Self {
            reader: Arc::new(ManualReader::builder().build()),
        }
    }
}

impl WasiMetricExporter {
    /// Exports metric data to a compatible host or component.
    pub fn export(&self) -> OTelSdkResult {
        let mut metrics = ResourceMetrics::default();
        // Scrape the metrics from the reader.
        self.reader.collect(&mut metrics)?;
        // Export to the host.
        wasi::otel::metrics::export(&metrics.into()).map_err(|e| e.into())
    }
}

impl MetricReader for WasiMetricExporter {
    /// Registers the metrics pipeline with this reader.
    ///
    /// Delegates to the embedded `ManualReader` to maintain shared ownership
    /// of metric data with an `SdkMeterProvider`.
    fn register_pipeline(&self, pipeline: std::sync::Weak<opentelemetry_sdk::metrics::Pipeline>) {
        self.reader.register_pipeline(pipeline)
    }

    /// Collects metrics from all registered instruments into the provided buffer.
    ///
    /// Delegates to the embedded `ManualReader` for the actual collection logic.
    fn collect(&self, rm: &mut ResourceMetrics) -> OTelSdkResult {
        self.reader.collect(rm)
    }

    /// This method is a no-op.
    fn force_flush(&self) -> OTelSdkResult {
        Ok(())
    }

    /// This method is a no-op.
    fn shutdown(&self) -> OTelSdkResult {
        Ok(())
    }

    /// This method is a no-op.
    fn shutdown_with_timeout(&self, _timeout: std::time::Duration) -> OTelSdkResult {
        Ok(())
    }

    /// Determines the temporality for a given instrument type.
    ///
    /// # Panics
    /// Panics if called with observable (async) instrument kinds, as these are not yet
    /// supported in WASI environments due to lack of background task execution.
    fn temporality(&self, kind: InstrumentKind) -> Temporality {
        match kind {
            InstrumentKind::ObservableCounter
            | InstrumentKind::ObservableGauge
            | InstrumentKind::ObservableUpDownCounter => {
                panic!("Async InstrumentKinds are not yet supported");
            }
            _ => self.reader.temporality(kind),
        }
    }
}

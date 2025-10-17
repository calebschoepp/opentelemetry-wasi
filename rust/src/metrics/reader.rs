use crate::wit::wasi::otel::metrics::collect;
use opentelemetry_sdk::{
    error::{OTelSdkError, OTelSdkResult},
    metrics::{
        data::ResourceMetrics, reader::MetricReader, InstrumentKind, ManualReader, Temporality,
    },
};
use std::sync::Arc;

#[derive(Debug)]
pub struct WasiMetricReader {
    reader: Arc<ManualReader>,
}

impl WasiMetricReader {
    pub fn new() -> Self {
        Self {
            reader: Arc::new(ManualReader::builder().build()),
        }
    }

    /// Creates a new instance sharing the same underlying data.
    ///
    /// This method is necessary because `SdkMeterProvider::builder().with_reader()`
    /// takes ownership of the reader, but we also need to retain access to call
    /// collection methods like `push()` on the same reader instance.
    pub fn reader(&self) -> Self {
        Self {
            reader: Arc::clone(&self.reader),
        }
    }

    /// Allows a compatible host or component to initiate the retrieval of metric data.
    pub fn pull() -> Result<ResourceMetrics, OTelSdkError> {
        unimplemented!()
    }

    /// Exports metric data to a compatible host or component.
    pub fn push(&self) -> OTelSdkResult {
        let mut metrics = ResourceMetrics::default();
        self.reader.collect(&mut metrics)?;
        // TODO: Maybe we rename the wasi `collect` method to `push` as well?
        collect(&metrics.into()).map_err(|e| e.into())
    }
}

impl MetricReader for WasiMetricReader {
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

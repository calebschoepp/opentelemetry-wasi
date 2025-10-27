use crate::wit::wasi;
use opentelemetry::otel_error;
use opentelemetry_sdk::{
    error::{OTelSdkError, OTelSdkResult},
    metrics::{
        data::ResourceMetrics, reader::MetricReader, InstrumentKind, ManualReader, Temporality,
    },
};
use std::sync::Arc;

/// A metric exporter that sends OpenTelemetry metrics to a WASI host.
///
/// This exporter wraps a `ManualReader` to use the reader's existing OpenTelemetry SDK internals
/// while providing a WASI export mechanism. The embedded reader handles metric
/// collection and the wrapper manages exports to the host.
///
/// # Default Example
/// ```ignore
/// let exporter = WasiMetricExporter::default();
/// let provider = SdkMeterProvider::builder().with_reader(exporter.clone());
/// // Measure something...
/// // Once the exporter is dropped, the metrics will automatically export to the host.
/// ```
///
/// # Manual Export Example
/// ```ignore
/// let exporter = WasiMetricExporter::builder()
///     .with_manual_export_only()
///     .build();
/// let provider = SdkMeterProvider::builder().with_reader(exporter.clone());
/// // Measure something...
/// exporter.export()?; // User must manually trigger export to host at some point before the end of the code.
/// ```
#[derive(Debug, Clone)]
pub struct WasiMetricExporter {
    reader: Arc<ManualReader>,
    export_on_drop: bool,
}

pub struct WasiMetricExporterBuilder {
    export_on_drop: bool,
}

impl WasiMetricExporterBuilder {
    pub fn new() -> Self {
        Self {
            export_on_drop: true,
        }
    }

    /// Configure the exporter to NOT automatically export when dropped.
    ///
    /// By default, the exporter will automatically export any collected metrics
    /// when it goes out of scope to prevent data loss. This method disables that
    /// behavior, requiring all exports to be triggered manually via [`export()`].
    ///
    /// Use this when you need precise control over when metrics are sent to the
    /// host.
    ///
    /// # Example
    /// ```ignore
    /// let exporter = WasiMetricExporter::builder()
    ///     .with_manual_export_only()
    ///     .build();
    /// // No export happens here
    /// drop(exporter);
    /// ```
    ///
    /// [`export()`]: WasiMetricExporter::export
    pub fn with_manual_export_only(mut self) -> Self {
        self.export_on_drop = false;
        self
    }

    /// Build the exporter.
    pub fn build(self) -> WasiMetricExporter {
        WasiMetricExporter {
            reader: Arc::new(ManualReader::builder().build()),
            export_on_drop: self.export_on_drop,
        }
    }
}

impl Default for WasiMetricExporter {
    fn default() -> Self {
        Self::builder().build()
    }
}

impl Drop for WasiMetricExporter {
    fn drop(&mut self) {
        if self.export_on_drop {
            _ = self.export();
        }
    }
}

impl WasiMetricExporter {
    /// Create a new builder for configuring a WasiMetricExporter.
    pub fn builder() -> WasiMetricExporterBuilder {
        WasiMetricExporterBuilder::new()
    }

    /// Exports metric data to a compatible host or component.
    pub fn export(&self) -> Result<(), OTelSdkError> {
        let mut metrics = ResourceMetrics::default();
        // Scrape the metrics from the reader.
        match self.reader.collect(&mut metrics) {
            Ok(_) => (),
            Err(sdk_error) => match sdk_error {
                OTelSdkError::AlreadyShutdown => {
                    otel_error!(name: "collect_already_shutdown", msg = "Shutdown has already been invoked.");
                    return Err(sdk_error);
                }
                OTelSdkError::Timeout(d) => {
                    otel_error!(name: "collect_timeout", msg = format!("Operation timed out after {} seconds.", d.as_secs()));
                    return Err(OTelSdkError::Timeout(d));
                }
                OTelSdkError::InternalFailure(e) => {
                    otel_error!(name: "collect_internal_failure", msg = format!("Operation failed due to an internal error: {}", e));
                    return Err(OTelSdkError::InternalFailure(e));
                }
            },
        }
        // Export to the host.
        match wasi::otel::metrics::export(&metrics.into()) {
            Ok(_) => Ok(()),
            Err(e) => {
                otel_error!(name: "export_internal_error", msg = format!("Operation failed due to an internal error: {}", e));
                return Err(OTelSdkError::InternalFailure(e));
            }
        }
    }
}

// Unless the `MetricReader` trait is specifically imported in the application code, these methods
// willl not be exposed to the end user. They are only meant to delegate to the manual reader or
// provide no-op methods that satisfy the trait requirements for an `SdkMeterProvider`.
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

    fn force_flush(&self) -> OTelSdkResult {
        Ok(())
    }

    fn shutdown(&self) -> OTelSdkResult {
        Ok(())
    }

    fn shutdown_with_timeout(&self, _timeout: std::time::Duration) -> OTelSdkResult {
        Ok(())
    }

    fn temporality(&self, kind: InstrumentKind) -> Temporality {
        self.reader.temporality(kind)
    }
}

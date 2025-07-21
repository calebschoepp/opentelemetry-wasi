use opentelemetry_sdk::metrics::{data::ResourceMetrics, exporter::PushMetricExporter, MetricResult, Temporality};
use crate::wit::wasi::otel::metrics::{export, force_flush, shutdown, temporality, MetricError, TemporalityT};
use async_trait::async_trait;
pub struct WasiExporter {}

#[async_trait]
impl PushMetricExporter for WasiExporter {
    async fn export(&self, metrics: &mut ResourceMetrics) -> MetricResult<()>{
        let converted = metrics.into();
        export(&converted).map_err(|e| match e {
            MetricError::Other(e) => opentelemetry_sdk::metrics::MetricError::Other(e),
            MetricError::Config(e) => opentelemetry_sdk::metrics::MetricError::Config(e),
            MetricError::InvalidInstrumentConfiguration(_) =>  {
                // TODO: Couldn't figure out how to get a String type into a static string, so I'm using a generic message for now
                opentelemetry_sdk::metrics::MetricError::InvalidInstrumentConfiguration("Invalid instrument configuration")
            }
        })
    }

    async fn force_flush(&self) -> MetricResult<()>{
        force_flush().map_err(|e| match e {
            MetricError::Other(e) => opentelemetry_sdk::metrics::MetricError::Other(e),
            MetricError::Config(e) => opentelemetry_sdk::metrics::MetricError::Config(e),
            MetricError::InvalidInstrumentConfiguration(_) =>  {
                // TODO: Couldn't figure out how to get a String type into a static string, so I'm using a generic message for now
                opentelemetry_sdk::metrics::MetricError::InvalidInstrumentConfiguration("Invalid instrument configuration")
            }
        })
    }

    fn shutdown(&self) -> MetricResult<()>{
        shutdown().map_err(|e| match e {
            MetricError::Other(e) => opentelemetry_sdk::metrics::MetricError::Other(e),
            MetricError::Config(e) => opentelemetry_sdk::metrics::MetricError::Config(e),
            MetricError::InvalidInstrumentConfiguration(_) =>  {
                // TODO: Couldn't figure out how to get a String type into a static string, so I'm using a generic message for now
                opentelemetry_sdk::metrics::MetricError::InvalidInstrumentConfiguration("Invalid instrument configuration")
            }
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
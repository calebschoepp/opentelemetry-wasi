use std::time::SystemTime;

use futures_core::future::BoxFuture;
use opentelemetry::trace::TraceError;
use opentelemetry_sdk::export::trace::{ExportResult, SpanData, SpanExporter};

use crate::wit::v2::observe::{self, Datetime};

#[derive(Debug)]
pub struct WasiExporter {}

impl WasiExporter {
    pub fn new() -> Self {
        Self {}
    }

    fn export_one(span_data: SpanData) -> anyhow::Result<()> {
        let start_since_the_epoch = span_data
            .start_time
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap();
        let end_since_the_epoch = span_data
            .end_time
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap();
        observe::emit_span(&observe::ReadOnlySpan {
            name: span_data.name.to_string(),
            span_context: observe::SpanContext {
                span_id: u64::from_be_bytes(span_data.span_context.span_id().to_bytes()),
                trace_id: span_data.span_context.trace_id().to_bytes().to_vec(),
            },
            span_kind: observe::SpanKind::Internal, // TODO
            start_time: Datetime {
                seconds: start_since_the_epoch.as_secs(),
                nanoseconds: start_since_the_epoch.subsec_nanos(),
            },
            end_time: Datetime {
                seconds: end_since_the_epoch.as_secs(),
                nanoseconds: end_since_the_epoch.subsec_nanos(),
            },
        });
        Ok(())
    }
}

impl SpanExporter for WasiExporter {
    fn export(&mut self, batch: Vec<SpanData>) -> BoxFuture<'static, ExportResult> {
        Box::pin(async move {
            for span_data in batch {
                if let Err(e) = Self::export_one(span_data) {
                    return Err(TraceError::Other(e.into()));
                }
            }
            Ok(())
        })
    }
}

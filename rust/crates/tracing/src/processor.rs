use std::sync::atomic::{AtomicBool, Ordering};

use opentelemetry_sdk::{error::OTelSdkResult, trace::SpanProcessor};

use crate::wit::wasi::otel::tracing::{on_end, on_start};

#[derive(Debug)]
pub struct WasiProcessor {
    is_shutdown: AtomicBool,
}

impl WasiProcessor {
    /// Create a new `WasiProcessor`.
    pub fn new() -> Self {
        Self {
            is_shutdown: AtomicBool::new(false),
        }
    }
}

impl Default for WasiProcessor {
    fn default() -> Self {
        Self::new()
    }
}

impl SpanProcessor for WasiProcessor {
    fn on_start(&self, span: &mut opentelemetry_sdk::trace::Span, _: &opentelemetry::Context) {
        println!("WasiProcessor::on_start invoked");
        if self.is_shutdown.load(Ordering::Relaxed) {
            return;
        }
        if let Some(span_data) = span.exported_data() {
            println!("SpanData is SOME");
            on_start(&span_data.span_context.into());
        } else {
            println!("SpanData is NONE");
        }
    }

    fn on_end(&self, span: opentelemetry_sdk::trace::SpanData) {
        if self.is_shutdown.load(Ordering::Relaxed) {
            return;
        }
        on_end(&span.into());
    }

    fn force_flush(&self) -> OTelSdkResult {
        if self.is_shutdown.load(Ordering::Relaxed) {
            return OTelSdkResult::Err(opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown);
        }
        Ok(())
    }

    fn shutdown(&self) -> OTelSdkResult {
        let result = self.force_flush();
        if self.is_shutdown.swap(true, Ordering::Relaxed) {
            return OTelSdkResult::Err(opentelemetry_sdk::error::OTelSdkError::AlreadyShutdown);
        }
        result
    }
}

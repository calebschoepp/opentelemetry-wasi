use std::sync::atomic::{AtomicBool, Ordering};

use opentelemetry::{
    otel_warn,
    trace::{TraceContextExt, TraceError},
};
use opentelemetry_sdk::trace::SpanProcessor;

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
    fn on_start(&self, span: &mut opentelemetry_sdk::trace::Span, cx: &opentelemetry::Context) {
        if self.is_shutdown.load(Ordering::Relaxed) {
            // this is a warning, as the user is trying to emit after the processor has been shutdown
            otel_warn!(
                name: "WasiProcessor.Emit.ProcessorShutdown",
            );
            return;
        }
        if let Some(span_data) = span.exported_data() {
            on_start(&span_data.into(), &cx.span().span_context().clone().into());
        }
    }

    fn on_end(&self, span: opentelemetry_sdk::export::trace::SpanData) {
        if self.is_shutdown.load(Ordering::Relaxed) {
            // this is a warning, as the user is trying to emit after the processor has been shutdown
            otel_warn!(
                name: "WasiProcessor.Emit.ProcessorShutdown",
            );
            return;
        }
        on_end(&span.into());
    }

    fn force_flush(&self) -> opentelemetry::trace::TraceResult<()> {
        if self.is_shutdown.load(Ordering::Relaxed) {
            return Err(TraceError::Other("Processor already shutdown".into()));
        }
        Ok(())
    }

    fn shutdown(&self) -> opentelemetry::trace::TraceResult<()> {
        let result = self.force_flush();
        if self.is_shutdown.swap(true, Ordering::Relaxed) {
            return Err(TraceError::Other("Processor already shutdown".into()));
        }
        result
    }
}

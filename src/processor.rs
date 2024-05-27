use std::sync::Mutex;

use opentelemetry::{
    global,
    trace::{TraceError, TraceResult},
    Context,
};
use opentelemetry_sdk::{
    export::trace::{SpanData, SpanExporter},
    trace::{Span, SpanProcessor},
};

#[derive(Debug)]
pub struct SimpleSpanProcessor {
    exporter: Mutex<Box<dyn SpanExporter>>,
}

impl SimpleSpanProcessor {
    pub fn new(exporter: Box<dyn SpanExporter>) -> Self {
        Self {
            exporter: Mutex::new(exporter),
        }
    }
}

impl SpanProcessor for SimpleSpanProcessor {
    fn on_start(&self, _span: &mut Span, _cx: &Context) {
        // Ignored
    }

    fn on_end(&self, span: SpanData) {
        if !span.span_context.is_sampled() {
            return;
        }

        let result = self
            .exporter
            .lock()
            .map_err(|_| TraceError::Other("SimpleSpanProcessor mutex poison".into()))
            .and_then(|mut exporter| futures_executor::block_on(exporter.export(vec![span])));

        if let Err(err) = result {
            global::handle_error(err);
        }
    }

    fn force_flush(&self) -> TraceResult<()> {
        // Nothing to flush for simple span processor.
        Ok(())
    }

    fn shutdown(&mut self) -> TraceResult<()> {
        if let Ok(mut exporter) = self.exporter.lock() {
            exporter.shutdown();
            Ok(())
        } else {
            Err(TraceError::Other(
                "SimpleSpanProcessor mutex poison at shutdown".into(),
            ))
        }
    }
}

use std::time::SystemTime;

use opentelemetry::global::ObjectSafeSpan;
use opentelemetry_sdk::trace::SpanProcessor;

use crate::wit::v2::observe;

#[derive(Debug)]
pub struct WasiProcessor {}

impl WasiProcessor {
    pub fn new() -> Self {
        Self {}
    }
}

impl Default for WasiProcessor {
    fn default() -> Self {
        Self::new()
    }
}

impl SpanProcessor for WasiProcessor {
    fn on_start(&self, span: &mut opentelemetry_sdk::trace::Span, _cx: &opentelemetry::Context) {
        println!(
            "Processor on_start passing span_id: {:?}",
            span.span_context().span_id().to_string()
        );
        observe::on_span_start(&observe::SpanContext {
            span_id: span.span_context().span_id().to_string(),
            trace_id: span.span_context().trace_id().to_string(),
            trace_flags: format!("{:x}", span.span_context().trace_flags()),
            is_remote: span.span_context().is_remote(),
            trace_state: "".to_string(), // TODO
        })
    }

    fn on_end(&self, span: opentelemetry_sdk::export::trace::SpanData) {
        println!("PROCESSOR ON_END FOR SPAN: {:?}\n\n", span.name.to_string());
        println!(
            "Processor on_end passing span_id: {:?}",
            span.span_context.span_id().to_string()
        );
        let start_since_the_epoch = span
            .start_time
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap();
        let end_since_the_epoch = span
            .end_time
            .duration_since(SystemTime::UNIX_EPOCH)
            .unwrap();
        observe::on_span_end(&observe::ReadOnlySpan {
            name: span.name.to_string(),
            span_context: observe::SpanContext {
                span_id: span.span_context.span_id().to_string(),
                trace_id: span.span_context.trace_id().to_string(),
                trace_flags: format!("{:x}", span.span_context.trace_flags()),
                is_remote: span.span_context.is_remote(),
                trace_state: "".to_string(), // TODO
            },
            parent_span_id: span.parent_span_id.to_string(),
            span_kind: observe::SpanKind::Internal, // TODO
            start_time: observe::Datetime {
                seconds: start_since_the_epoch.as_secs(),
                nanoseconds: start_since_the_epoch.subsec_nanos(),
            },
            end_time: observe::Datetime {
                seconds: end_since_the_epoch.as_secs(),
                nanoseconds: end_since_the_epoch.subsec_nanos(),
            },
            attributes: vec![],
            otel_resource: observe::OtelResource {
                attrs: vec![],
                schema_url: None,
            },
        })
    }

    fn force_flush(&self) -> opentelemetry::trace::TraceResult<()> {
        Ok(())
    }

    fn shutdown(&mut self) -> opentelemetry::trace::TraceResult<()> {
        Ok(())
    }
}

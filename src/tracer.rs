use std::sync::Arc;

use opentelemetry::{
    trace::{SpanContext, TraceContextExt},
    InstrumentationLibrary,
};
use tracing_opentelemetry::PreSampledTracer;

use crate::span::Span;

pub struct Tracer {
    _instrumentation_lib: Arc<InstrumentationLibrary>,
}

impl Tracer {
    pub(crate) fn new(instrumentation_lib: Arc<InstrumentationLibrary>) -> Self {
        Self {
            _instrumentation_lib: instrumentation_lib,
        }
    }
}

impl opentelemetry::trace::Tracer for Tracer {
    type Span = Span;

    fn build_with_context(
        &self,
        builder: opentelemetry::trace::SpanBuilder,
        parent_cx: &opentelemetry::Context,
    ) -> Self::Span {
        Span::new(builder.name, parent_cx.span().span_context().clone())
    }
}

impl PreSampledTracer for Tracer {
    // TODO: Incorrectly always assume we're sampling here
    fn sampled_context(
        &self,
        data: &mut tracing_opentelemetry::OtelData,
    ) -> opentelemetry::Context {
        data.parent_cx.clone()
    }

    fn new_trace_id(&self) -> opentelemetry::trace::TraceId {
        opentelemetry::trace::TraceId::INVALID
    }

    fn new_span_id(&self) -> opentelemetry::trace::SpanId {
        opentelemetry::trace::SpanId::INVALID
    }
}

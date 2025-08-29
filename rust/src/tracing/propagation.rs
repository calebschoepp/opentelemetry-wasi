use crate::wit::wasi::otel::tracing::outer_span_context;
use opentelemetry::{trace::TraceContextExt, Context};

pub trait WasiPropagator {
    fn extract(&self, cx: &Context) -> Context;
}

pub struct TraceContextPropagator {}

impl TraceContextPropagator {
    pub fn new() -> Self {
        Self {}
    }
}

impl Default for TraceContextPropagator {
    fn default() -> Self {
        Self::new()
    }
}

impl WasiPropagator for TraceContextPropagator {
    fn extract(&self, cx: &Context) -> Context {
        cx.with_remote_span_context(outer_span_context().into())
    }
}

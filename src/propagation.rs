use crate::wit::v2::observe::get_parent_span_context;
use opentelemetry::{
    trace::{
        SpanContext as OtelSpanContext, SpanId, TraceContextExt, TraceFlags, TraceId, TraceState,
    },
    Context, ContextGuard,
};
use std::str::FromStr;

pub fn extract_trace_context() -> anyhow::Result<ContextGuard> {
    let other_sc = get_parent_span_context();

    let sc = OtelSpanContext::new(
        TraceId::from_hex(&other_sc.trace_id)?,
        SpanId::from_hex(&other_sc.span_id)?,
        TraceFlags::SAMPLED,
        other_sc.is_remote,
        TraceState::from_str(&other_sc.trace_state)?,
    );

    Ok(Context::current().with_remote_span_context(sc).attach())
}

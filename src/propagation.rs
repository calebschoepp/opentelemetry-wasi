use opentelemetry::{
    trace::{SpanContext as OtelSpanContext, TraceContextExt, TraceFlags, TraceState},
    Context, ContextGuard,
};

use crate::wit::v2::observe::get_parent;

pub fn extract_trace_context() -> ContextGuard {
    let other_sc = get_parent();

    let trace_id_array: [u8; 16] = other_sc
        .trace_id
        .into_iter()
        .collect::<Vec<u8>>()
        .try_into()
        .unwrap();

    let sc = OtelSpanContext::new(
        u128::from_be_bytes(trace_id_array).into(),
        other_sc.span_id.into(),
        TraceFlags::SAMPLED,
        false, // TODO: Is this correct?
        TraceState::default(),
    );

    Context::current().with_remote_span_context(sc).attach()
}

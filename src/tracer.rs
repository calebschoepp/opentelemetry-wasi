// use std::sync::Arc;

// use opentelemetry::InstrumentationLibrary;

// use crate::span::Span;

// pub struct Tracer {
//     _instrumentation_lib: Arc<InstrumentationLibrary>,
// }

// impl Tracer {
//     pub(crate) fn new(instrumentation_lib: Arc<InstrumentationLibrary>) -> Self {
//         Self {
//             _instrumentation_lib: instrumentation_lib,
//         }
//     }
// }

// impl opentelemetry::trace::Tracer for Tracer {
//     type Span = Span;

//     fn build_with_context(
//         &self,
//         builder: opentelemetry::trace::SpanBuilder,
//         _parent_cx: &opentelemetry::Context,
//     ) -> Self::Span {
//         Span::new(builder.name)
//     }
// }

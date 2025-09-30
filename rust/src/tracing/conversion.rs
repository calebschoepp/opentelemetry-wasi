use crate::wit::wasi::otel::tracing::*;

impl From<opentelemetry_sdk::trace::SpanData> for SpanData {
    fn from(value: opentelemetry_sdk::trace::SpanData) -> Self {
        Self {
            span_context: value.span_context.into(),
            parent_span_id: value.parent_span_id.to_string(),
            span_kind: value.span_kind.into(),
            name: value.name.to_string(),
            start_time: value.start_time.into(),
            end_time: value.end_time.into(),
            attributes: value.attributes.into_iter().map(Into::into).collect(),
            events: value.events.events.into_iter().map(Into::into).collect(),
            links: value.links.links.into_iter().map(Into::into).collect(),
            status: value.status.into(),
            instrumentation_scope: value.instrumentation_scope.into(),
            dropped_attributes: value.dropped_attributes_count,
            dropped_events: value.events.dropped_count,
            dropped_links: value.links.dropped_count,
        }
    }
}

impl From<opentelemetry::trace::SpanContext> for SpanContext {
    fn from(value: opentelemetry::trace::SpanContext) -> Self {
        Self {
            trace_id: format!("{:x}", value.trace_id()),
            span_id: format!("{:x}", value.span_id()),
            trace_flags: value.trace_flags().into(),
            is_remote: value.is_remote(),
            trace_state: value
                .trace_state()
                .header()
                .split(',')
                .filter_map(|s| {
                    if let Some((key, value)) = s.split_once('=') {
                        Some((key.to_string(), value.to_string()))
                    } else {
                        None
                    }
                })
                .collect(),
        }
    }
}

impl From<SpanContext> for opentelemetry::trace::SpanContext {
    fn from(value: SpanContext) -> Self {
        let trace_id = opentelemetry::trace::TraceId::from_hex(&value.trace_id)
            .unwrap_or(opentelemetry::trace::TraceId::INVALID);
        let span_id = opentelemetry::trace::SpanId::from_hex(&value.span_id)
            .unwrap_or(opentelemetry::trace::SpanId::INVALID);
        let trace_state = opentelemetry::trace::TraceState::from_key_value(value.trace_state)
            .unwrap_or_else(|_| opentelemetry::trace::TraceState::default());
        Self::new(
            trace_id,
            span_id,
            value.trace_flags.into(),
            value.is_remote,
            trace_state,
        )
    }
}

impl From<opentelemetry::trace::TraceFlags> for TraceFlags {
    fn from(value: opentelemetry::trace::TraceFlags) -> Self {
        if value.is_sampled() {
            TraceFlags::SAMPLED
        } else {
            TraceFlags::empty()
        }
    }
}

impl From<TraceFlags> for opentelemetry::trace::TraceFlags {
    fn from(value: TraceFlags) -> Self {
        Self::new(value.bits())
    }
}

impl From<opentelemetry::trace::SpanKind> for SpanKind {
    fn from(value: opentelemetry::trace::SpanKind) -> Self {
        match value {
            opentelemetry::trace::SpanKind::Client => Self::Client,
            opentelemetry::trace::SpanKind::Server => Self::Server,
            opentelemetry::trace::SpanKind::Producer => Self::Producer,
            opentelemetry::trace::SpanKind::Consumer => Self::Consumer,
            opentelemetry::trace::SpanKind::Internal => Self::Internal,
        }
    }
}

impl From<opentelemetry::trace::Event> for Event {
    fn from(value: opentelemetry::trace::Event) -> Self {
        Self {
            name: value.name.to_string(),
            time: value.timestamp.into(),
            attributes: value.attributes.into_iter().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry::trace::Link> for Link {
    fn from(value: opentelemetry::trace::Link) -> Self {
        Self {
            span_context: value.span_context.into(),
            attributes: value.attributes.into_iter().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry::trace::Status> for Status {
    fn from(value: opentelemetry::trace::Status) -> Self {
        match value {
            opentelemetry::trace::Status::Unset => Self::Unset,
            opentelemetry::trace::Status::Error { description } => {
                Self::Error(description.to_string())
            }
            opentelemetry::trace::Status::Ok => Self::Ok,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use core::panic;
    use opentelemetry::trace as O;
    use opentelemetry_sdk::trace as Sdk;
    use std::time;

    #[test]
    fn span_data_conversion() {
        // TODO: there were some fields that I didn't test due to using the defaults:
        // - events
        // - instrumentation_scope.{attributes, schema_url, version}
        // - links
        // - span_context.trace_state
        let otel_span = Sdk::SpanData {
            attributes: vec![opentelemetry::KeyValue::new("foo", "bar")],
            dropped_attributes_count: 6,
            end_time: time::UNIX_EPOCH + time::Duration::from_secs(1234567890),
            events: Sdk::SpanEvents::default(),
            instrumentation_scope: opentelemetry::InstrumentationScope::builder("tst-scp").build(),
            links: Sdk::SpanLinks::default(),
            name: std::borrow::Cow::Borrowed("test-name"),
            parent_span_id: opentelemetry::SpanId::from_bytes([2u8; 8]),
            span_context: O::SpanContext::new(
                O::TraceId::from_bytes([0u8; 16]),
                O::SpanId::from_bytes([1u8; 8]),
                O::TraceFlags::SAMPLED,
                true,
                O::TraceState::default(),
            ),
            span_kind: O::SpanKind::Internal,
            start_time: time::UNIX_EPOCH + time::Duration::from_secs(9876543210),
            status: O::Status::Ok,
        };

        let span: SpanData = otel_span.into();

        // Root
        assert_eq!(span.dropped_attributes, 6u32);
        assert_eq!(span.end_time.seconds, 1234567890u64);
        assert_eq!(span.name, "test-name".to_string());
        assert_eq!(span.parent_span_id, "0202020202020202");
        assert_eq!(span.span_kind, SpanKind::Internal);
        assert_eq!(span.start_time.seconds, 9876543210u64);
        assert!(matches!(span.status, Status::Ok));

        // SpanContext
        assert_eq!(span.span_context.trace_id, "0");
        assert_eq!(span.span_context.span_id, "101010101010101");
        assert_eq!(span.span_context.trace_flags, TraceFlags::SAMPLED);
        assert!(span.span_context.is_remote);

        // Attributes
        assert!(span.attributes.len() == 1);
        assert_eq!(span.attributes[0].key, "foo".to_string());
        match &span.attributes[0].value {
            crate::wit::wasi::otel::types::Value::String(s) => assert_eq!(s, "bar"),
            _ => panic!("Expected Value::String, got {:?}", span.attributes[0].value),
        }

        // InstrumentationScope
        assert_eq!(span.instrumentation_scope.name, "tst-scp".to_string());
    }
}

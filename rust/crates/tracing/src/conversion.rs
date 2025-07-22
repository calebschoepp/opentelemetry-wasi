use std::time::UNIX_EPOCH;

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

impl From<std::time::SystemTime> for Datetime {
    fn from(value: std::time::SystemTime) -> Self {
        let duration_since_epoch = value
            .duration_since(UNIX_EPOCH)
            .expect("SystemTime should be after UNIX EPOCH");
        Self {
            seconds: duration_since_epoch.as_secs(),
            nanoseconds: duration_since_epoch.subsec_nanos(),
        }
    }
}

impl From<opentelemetry::KeyValue> for KeyValue {
    fn from(value: opentelemetry::KeyValue) -> Self {
        Self {
            key: value.key.to_string(),
            value: value.value.into(),
        }
    }
}

impl From<&opentelemetry::KeyValue> for KeyValue {
    fn from(value: &opentelemetry::KeyValue) -> Self {
        Self {
            key: value.key.to_string(),
            value: value.value.clone().into(),
        }
    }
}

impl From<opentelemetry::Value> for Value {
    fn from(value: opentelemetry::Value) -> Self {
        match value {
            opentelemetry::Value::Bool(v) => Self::Bool(v),
            opentelemetry::Value::I64(v) => Self::S64(v),
            opentelemetry::Value::F64(v) => Self::F64(v),
            opentelemetry::Value::String(v) => Self::String(v.to_string()),
            opentelemetry::Value::Array(v) => match v {
                opentelemetry::Array::Bool(items) => Self::BoolArray(items),
                opentelemetry::Array::I64(items) => Self::S64Array(items),
                opentelemetry::Array::F64(items) => Self::F64Array(items),
                opentelemetry::Array::String(items) => {
                    Self::StringArray(items.into_iter().map(Into::into).collect())
                }
                _ => unimplemented!(),
            },
            _ => unimplemented!(),
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

impl From<opentelemetry::InstrumentationScope> for InstrumentationScope {
    fn from(value: opentelemetry::InstrumentationScope) -> Self {
        Self {
            name: value.name().to_string(),
            version: value.version().map(Into::into),
            schema_url: value.schema_url().map(Into::into),
            attributes: value.attributes().map(Into::into).collect(),
        }
    }
}

use crate::wit::wasi::{
    clocks::wall_clock::Datetime,
    otel::types::{InstrumentationScope, KeyValue, Value},
};
use std::time::UNIX_EPOCH;

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

impl From<(&opentelemetry::Key, &opentelemetry::Value)> for KeyValue {
    fn from(value: (&opentelemetry::Key, &opentelemetry::Value)) -> Self {
        Self {
            key: value.0.to_string(),
            value: value.1.into(),
        }
    }
}

impl From<&opentelemetry::Value> for Value {
    fn from(value: &opentelemetry::Value) -> Self {
        match value {
            opentelemetry::Value::Bool(v) => Self::Bool(v.to_owned()),
            opentelemetry::Value::I64(v) => Self::S64(v.to_owned()),
            opentelemetry::Value::F64(v) => Self::F64(v.to_owned()),
            opentelemetry::Value::String(v) => Self::String(v.to_string()),
            opentelemetry::Value::Array(v) => match v {
                opentelemetry::Array::Bool(items) => Self::BoolArray(items.to_owned()),
                opentelemetry::Array::I64(items) => Self::S64Array(items.to_owned()),
                opentelemetry::Array::F64(items) => Self::F64Array(items.to_owned()),
                opentelemetry::Array::String(items) => {
                    Self::StringArray(items.iter().map(|e| e.to_string()).collect())
                }
                _ => unimplemented!(),
            },
            _ => unimplemented!(),
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

impl From<&opentelemetry::InstrumentationScope> for InstrumentationScope {
    fn from(value: &opentelemetry::InstrumentationScope) -> Self {
        Self {
            name: value.name().to_string(),
            version: value.version().map(Into::into),
            schema_url: value.schema_url().map(Into::into),
            attributes: value.attributes().map(Into::into).collect(),
        }
    }
}

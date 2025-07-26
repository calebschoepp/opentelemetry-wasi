use crate::wit::wasi::otel::metrics::*;
use crate::wit::wasi::otel::types::{KeyValue, Value};
use std::time::UNIX_EPOCH;

impl From<&mut opentelemetry_sdk::metrics::data::ResourceMetrics> for ResourceMetrics {
    fn from(value: &mut opentelemetry_sdk::metrics::data::ResourceMetrics) -> Self {
        Self {
            resource: value.resource.to_owned().into(),
            scope_metrics: value.scope_metrics
            .iter()
            .to_owned()
            .map(Into::into)
            .collect(),
        }
    }
}

impl From<opentelemetry_sdk::resource::Resource> for Resource {
    fn from(value: opentelemetry_sdk::resource::Resource) -> Self {
        Self {
            inner: ResourceInner {
                attributes: value.into_iter().map(Into::into).collect(),
                schema_url: match value.schema_url() {
                    Some(v ) => Some(v.to_string()),
                    None => None,
                },
            },
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

impl From<&opentelemetry::KeyValue> for KeyValue {
    fn from(value: &opentelemetry::KeyValue) -> Self {
        Self {
            key: value.key.to_owned().into(),
            value: value.value.to_owned().into(),
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
                    Self::StringArray(items.to_owned().into_iter().map(Into::into).collect())
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

impl From<opentelemetry_sdk::metrics::data::ScopeMetrics> for ScopeMetrics {
    fn from(value: opentelemetry_sdk::metrics::data::ScopeMetrics) -> Self {
        Self {
            scope: value.scope.into(),
            metrics: value.metrics.into_iter().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::ScopeMetrics> for ScopeMetrics {
    fn from(value: &opentelemetry_sdk::metrics::data::ScopeMetrics) -> Self {
        Self {
            scope: value.scope.to_owned().into(),
            metrics: value.metrics.iter().to_owned().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry::InstrumentationScope> for InstrumentationScope {
    fn from(value: opentelemetry::InstrumentationScope) -> Self {
        Self {
            name: value.name().to_string(),
            version: value.name().to_string(),
            schema_url: value.name().to_string(),
            attributes: value.attributes().into_iter().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry_sdk::metrics::data::Metric> for Metric {
    fn from(value: opentelemetry_sdk::metrics::data::Metric) -> Self {
        Self {
            name: value.name.to_string(),
            description: value.description.to_string(),
            unit: value.unit.to_string(),
            data: value.data.into(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Metric> for Metric {
    fn from(value: &opentelemetry_sdk::metrics::data::Metric) -> Self {
        Self {
            name: value.name.to_owned().to_string(),
            description: value.description.to_owned().to_string(),
            unit: value.unit.to_owned().to_string(),
            data: value.data.into(),
        }
    }
}

impl From<opentelemetry_sdk::metrics::data::Aggregation> for MetricData {
    fn from(value: opentelemetry_sdk::metrics::data::Aggregation) -> Self {
       todo!()
    }
}
use crate::wit::wasi::otel::metrics::*;

impl From<opentelemetry_sdk::metrics::data::ResourceMetrics> for ResourceMetrics {
    fn from(value: opentelemetry_sdk::metrics::data::ResourceMetrics) -> Self {
        Self {
            resource: value.resource.into(),
            scope_metrics: value.scope_metrics.into_iter().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry_sdk::resource::Resource> for Resource {
    fn from(value: opentelemetry_sdk::resource::Resource) -> Self {
        let mut attrs: Vec<KeyValue> = Vec::new();
        value.into_iter().map(|v| {
            attrs.push(KeyValue { key: v.0.to_string(), value: v.1.into() })
        });

        Self {
            inner: ResourceInner {
                attrs: attrs,
                schema_url: value.schema_url().map(Into::into),
            }
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

impl From<opentelemetry_sdk::metrics::data::Metric> for Metric {
    fn from(value: opentelemetry_sdk::metrics::data::Metric) -> Self {
        Self {
            name: value.name.into(),
            description: value.description.into(),
            unit: value.unit.into(),
            data: value.data.into(),
        }
    }
}

// TODO: duplicated code
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

// TODO: duplicated code
impl From<opentelemetry::KeyValue> for KeyValue {
    fn from(value: opentelemetry::KeyValue) -> Self {
        Self {
            key: value.key.to_string(),
            value: value.value.into(),
        }
    }
}

// TODO: duplicated code
impl From<&opentelemetry::KeyValue> for KeyValue {
    fn from(value: &opentelemetry::KeyValue) -> Self {
        Self {
            key: value.key.to_string(),
            value: value.value.clone().into(),
        }
    }
}

// TODO: duplicated code
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
use crate::wit::wasi::{
    clocks::wall_clock::Datetime,
    otel::types::{InstrumentationScope, KeyValue, Resource},
};
use serde::{ser::SerializeSeq, Serialize};
use std::time::UNIX_EPOCH;

/// Converts a Serde-serializable type to JSON.
pub fn to_json<T: Serialize>(v: &T) -> String {
    serde_json::to_string(v).unwrap_or_else(|e| {
        panic!(
            "failed to serialize {} to JSON: {}",
            std::any::type_name::<T>(),
            e
        )
    })
}

#[macro_export]
#[doc(hidden)]
/// Serializes an array using Serde.
macro_rules! serialize_seq {
    ($array:expr, $serializer:expr, |$elem:ident| $transform:expr) => {{
        let mut seq = $serializer.serialize_seq(Some($array.len()))?;
        for $elem in $array.iter() {
            // This is for elements that require additional transformations
            // (e.g., call the `as_str()` method on the OpenTelemetry `StringValue` type).
            seq.serialize_element($transform)?;
        }
        seq.end()
    }};
    ($array:expr, $serializer:expr) => {{
        let mut seq = $serializer.serialize_seq(Some($array.len()))?;
        for elem in $array.iter() {
            seq.serialize_element(elem)?;
        }
        seq.end()
    }};
}

struct ValueWrapper<'a>(&'a opentelemetry::Value);
impl<'a> Serialize for ValueWrapper<'a> {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        match self.0 {
            opentelemetry::Value::Bool(v) => serializer.serialize_bool(*v),
            opentelemetry::Value::F64(v) => serializer.serialize_f64(*v),
            opentelemetry::Value::I64(v) => serializer.serialize_i64(*v),
            opentelemetry::Value::String(v) => serializer.serialize_str(v.as_str()),
            opentelemetry::Value::Array(arr) => match arr {
                opentelemetry::Array::Bool(v) => serialize_seq!(v, serializer),
                opentelemetry::Array::F64(v) => serialize_seq!(v, serializer),
                opentelemetry::Array::I64(v) => serialize_seq!(v, serializer),
                opentelemetry::Array::String(v) => {
                    serialize_seq!(v, serializer, |sv| sv.as_str())
                }
                _ => unimplemented!(),
            },
            _ => unimplemented!(),
        }
    }
}

impl From<&opentelemetry::KeyValue> for KeyValue {
    fn from(value: &opentelemetry::KeyValue) -> Self {
        Self {
            key: value.key.to_string(),
            value: to_json(&ValueWrapper(&value.value)),
        }
    }
}

impl From<(&opentelemetry::Key, &opentelemetry::Value)> for KeyValue {
    fn from(value: (&opentelemetry::Key, &opentelemetry::Value)) -> Self {
        Self {
            key: value.0.to_string(),
            value: to_json(&ValueWrapper(value.1)),
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

impl From<&opentelemetry_sdk::resource::Resource> for Resource {
    fn from(value: &opentelemetry_sdk::resource::Resource) -> Self {
        Self {
            attributes: value.into_iter().map(Into::into).collect(),
            schema_url: value.schema_url().map(Into::into),
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    // Converts an opentelemetry Value type to a JSON string.
    macro_rules! otel_value_to_serde_string {
        ($v:expr) => {
            serde_json::to_string(&ValueWrapper(&$v)).unwrap()
        };
    }

    #[test]
    fn seserialize_otel_value_to_string() {
        use opentelemetry::{Array, Value};

        // Test bool
        let boolean = Value::Bool(false);
        assert_eq!("false", otel_value_to_serde_string!(boolean));
        let bool_arr = Value::Array(Array::Bool(vec![false, true, true]));
        assert_eq!("[false,true,true]", otel_value_to_serde_string!(bool_arr));

        // Test i64
        let int = Value::I64(6);
        assert_eq!("6", otel_value_to_serde_string!(int));
        let int_arr = Value::Array(Array::I64(vec![1, 2, 3, 4]));
        assert_eq!("[1,2,3,4]", otel_value_to_serde_string!(int_arr));

        // Test f64
        let float = Value::F64(123.456);
        assert_eq!("123.456", otel_value_to_serde_string!(float));
        let float_arr = Value::Array(Array::F64(vec![1.0, 2.1, 3.2, 4.3]));
        assert_eq!("[1.0,2.1,3.2,4.3]", otel_value_to_serde_string!(float_arr));

        // Test String
        let str = Value::String("Test".into());
        assert_eq!("\"Test\"", otel_value_to_serde_string!(str));
        let str_arr = Value::Array(opentelemetry::Array::String(vec![
            "Hello, world!".into(),
            "Goodnight, moon.".into(),
        ]));
        assert_eq!(
            "[\"Hello, world!\",\"Goodnight, moon.\"]",
            otel_value_to_serde_string!(str_arr)
        );
    }
}

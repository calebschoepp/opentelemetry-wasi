use crate::{serialize_seq, types::to_json, wit::wasi::otel::logs::*};
use base64::Engine;
use serde::{
    ser::{SerializeMap, SerializeSeq},
    Serialize,
};

pub fn to_wasi_log_record(
    record: &opentelemetry_sdk::logs::SdkLogRecord,
    scope: &opentelemetry::InstrumentationScope,
    resource: Option<&opentelemetry_sdk::Resource>,
) -> LogRecord {
    let (trace_id, span_id, trace_flags) = record
        .trace_context()
        .as_ref()
        .map(|tc| {
            (
                format!("{:x}", tc.trace_id).into(),
                format!("{:x}", tc.span_id).into(),
                tc.trace_flags.map(Into::into),
            )
        })
        .unwrap_or((None, None, None));

    LogRecord {
        timestamp: record.timestamp().map(Into::into),
        body: record.body().map(|e| to_json(&AnyValueWrapper(e))),
        event_name: record.event_name().map(|e| e.to_string()),
        observed_timestamp: record.observed_timestamp().map(Into::into),
        severity_text: record.severity_text().map(|e| e.to_string()),
        severity_number: record.severity_number().map(|e| e as u8),
        attributes: record.attributes_iter().map(|e| Some(e.into())).collect(),
        instrumentation_scope: Some(scope.into()),
        resource: resource.map(Into::into),
        trace_id,
        span_id,
        trace_flags,
    }
}

impl From<&(opentelemetry::Key, opentelemetry::logs::AnyValue)> for KeyValue {
    fn from(value: &(opentelemetry::Key, opentelemetry::logs::AnyValue)) -> Self {
        Self {
            key: value.0.to_string(),
            value: to_json(&AnyValueWrapper(&value.1)),
        }
    }
}

struct AnyValueWrapper<'a>(&'a opentelemetry::logs::AnyValue);
impl<'a> Serialize for AnyValueWrapper<'a> {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        match self.0 {
            opentelemetry::logs::AnyValue::Boolean(v) => serializer.serialize_bool(*v),
            opentelemetry::logs::AnyValue::Int(v) => serializer.serialize_i64(*v),
            opentelemetry::logs::AnyValue::Double(v) => serializer.serialize_f64(*v),
            opentelemetry::logs::AnyValue::String(v) => serializer.serialize_str(v.as_str()),
            opentelemetry::logs::AnyValue::Bytes(bytes) => {
                // This is a workaround for JSON not having a way to differentiate between an array of bytes and an array of integers.
                let encoded = base64::engine::general_purpose::STANDARD.encode(bytes.as_ref());
                serializer
                    .serialize_str(&format!("data:application/octet-stream;base64,{}", encoded))
            }
            opentelemetry::logs::AnyValue::ListAny(list) => {
                serialize_seq!(list, serializer, |v| &AnyValueWrapper(v))
            }
            opentelemetry::logs::AnyValue::Map(map) => {
                let mut result_map = serializer.serialize_map(Some(map.len()))?;
                for (k, v) in map.iter() {
                    result_map.serialize_entry(k.as_str(), &AnyValueWrapper(v))?;
                }
                result_map.end()
            }
            _ => unimplemented!(),
        }
    }
}

#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    use super::*;
    use opentelemetry::{logs::AnyValue, Key};

    #[test]
    /// Tests that the AnyValue serializes into a string correctly.
    fn serialize_otel_log_any_value_to_string() {
        let mut hm: Box<HashMap<Key, AnyValue>> = Box::default();

        hm.insert(Key::new("key1"), AnyValue::Boolean(false));
        hm.insert(Key::new("key2"), AnyValue::Double(123.456));
        hm.insert(Key::new("key3"), AnyValue::Int(41));
        hm.insert(
            Key::new("key4"),
            AnyValue::Bytes(Box::new(b"Hello, world!".to_vec())),
        );
        hm.insert(
            Key::new("key5"),
            AnyValue::String("This is a string".into()),
        );
        hm.insert(
            Key::new("key6"),
            AnyValue::ListAny(Box::new(vec![
                AnyValue::Int(1),
                AnyValue::Int(2),
                AnyValue::Int(3),
            ])),
        );

        let mut nested_hm: Box<HashMap<Key, AnyValue>> = Box::default();

        nested_hm.insert(
            Key::new("nestedkey1"),
            AnyValue::String("Hello, from within!".into()),
        );

        hm.insert(Key::new("key7"), AnyValue::Map(nested_hm));

        let json_str =
            serde_json::to_string(&AnyValueWrapper(&opentelemetry::logs::AnyValue::Map(hm)))
                .unwrap();

        let actual: serde_json::Value = serde_json::from_str(&json_str).unwrap();
        let expected: serde_json::Value = serde_json::json!({
            "key1": false,
            "key2": 123.456,
            "key3": 41,
            //'Hello, world!' encoded to base64
            "key4": "data:application/octet-stream;base64,SGVsbG8sIHdvcmxkIQ==",
            "key5": "This is a string",
            "key6": [1, 2, 3],
            "key7": {
                "nestedkey1": "Hello, from within!"
            }
        });

        assert_eq!(actual, expected);
    }
}

use core::panic;
use std::{collections::HashMap, time::UNIX_EPOCH};

use crate::wit::wasi::otel::logs::*;

impl From<&mut opentelemetry_sdk::logs::SdkLogRecord> for LogRecord {
    fn from(value: &mut opentelemetry_sdk::logs::SdkLogRecord) -> Self {
        Self {
            event_name: match value.event_name() {
                Some(v) => Some(v.to_string()),
                None => None,
            },
            timestamp: match value.timestamp() {
                Some(v) => Some(v.into()),
                None => None,
            },
            observed_timestamp: match value.observed_timestamp() {
                Some(v) => Some(v.into()),
                None => None,
            },
            severity: match value.severity_number() {
                Some(v) => Some(v.into()),
                None => None,
            },
            severity_text: match value.severity_text() {
                Some(v) => Some(v.to_string()),
                None => None,
            },
            body: match value.body() {
                Some(v) => Some(v.into()),
                None => None,
            },
            attributes: value.attributes_iter().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry::logs::Severity> for Severity {
    fn from(value: opentelemetry::logs::Severity) -> Self {
        use opentelemetry::logs::Severity as S;
        match value {
            S::Trace => Severity::Trace,
            S::Trace2 => Severity::Trace2,
            S::Trace3 => Severity::Trace3,
            S::Trace4 => Severity::Trace4,
            S::Debug => Severity::Debug,
            S::Debug2 => Severity::Debug2,
            S::Debug3 => Severity::Debug3,
            S::Debug4 => Severity::Debug4,
            S::Info => Severity::Info,
            S::Info2 => Severity::Info2,
            S::Info3 => Severity::Info3,
            S::Info4 => Severity::Info4,
            S::Warn => Severity::Warn,
            S::Warn2 => Severity::Warn2,
            S::Warn3 => Severity::Warn3,
            S::Warn4 => Severity::Warn4,
            S::Error => Severity::Error,
            S::Error2 => Severity::Error2,
            S::Error3 => Severity::Error3,
            S::Error4 => Severity::Error4,
            S::Fatal => Severity::Fatal,
            S::Fatal2 => Severity::Fatal2,
            S::Fatal3 => Severity::Fatal3,
            S::Fatal4 => Severity::Fatal4,
        }
    }
}

impl From<opentelemetry::logs::AnyValue> for LogAny {
    fn from(value: opentelemetry::logs::AnyValue) -> Self {
        use opentelemetry::logs::AnyValue as A;
        match value {
            A::Int(v) => LogAny::Value(v.into()),
            A::Double(v) => LogAny::Value(v.into()),
            A::String(v) => LogAny::Value(v.into()),
            A::Boolean(v) => LogAny::Value(v.into()),
            A::Bytes(v) => LogAny::List(v.into()),
            A::ListAny(v) => LogAny::List(v.into()),
            A::Map(v) => LogAny::Map(map_to_kv_list(v)),
            _ => panic!("unsupported data type"),
        }
    }
}

impl From<&opentelemetry::logs::AnyValue> for LogAny {
    fn from(value: &opentelemetry::logs::AnyValue) -> Self {
        use opentelemetry::logs::AnyValue as A;
        match value {
            A::Int(v) => LogAny::Value(v.to_owned().into()),
            A::Double(v) => LogAny::Value(v.to_owned().into()),
            A::String(v) => LogAny::Value(v.to_owned().into()),
            A::Boolean(v) => LogAny::Value(v.to_owned().into()),
            A::Bytes(v) => LogAny::List(v.to_owned().into()),
            A::ListAny(v) => LogAny::List(v.to_owned().into()),
            A::Map(v) => LogAny::Map(map_to_kv_list(v.to_owned())),
            _ => panic!("unsupported data type"),
        }
    }
}

impl From<i64> for LogValue {
    fn from(value: i64) -> Self {
        Self::Int(value)
    }
}

impl From<f64> for LogValue {
    fn from(value: f64) -> Self {
        Self::Double(value)
    }
}

impl From<opentelemetry::StringValue> for LogValue {
    fn from(value: opentelemetry::StringValue) -> Self {
        Self::String(value.to_string())
    }
}

impl From<bool> for LogValue {
    fn from(value: bool) -> Self {
        Self::Boolean(value)
    }
}

impl From<Box<Vec<u8>>> for LogList {
    fn from(value: Box<Vec<u8>>) -> Self {
        Self::Bytes(*value)
    }
}

impl From<Box<Vec<opentelemetry::logs::AnyValue>>> for LogList {
    fn from(value: Box<Vec<opentelemetry::logs::AnyValue>>) -> Self {
        LogList::List(value.into_iter().map(Into::into).collect())
    }
}

impl From<opentelemetry::logs::AnyValue> for LogValue {
    fn from(value: opentelemetry::logs::AnyValue) -> Self {
        use opentelemetry::logs::AnyValue as A;
        match value {
            A::Int(v) => LogValue::Int(v),
            A::Double(v) => LogValue::Double(v),
            A::String(v) => LogValue::String(v.to_string()),
            A::Boolean(v) => LogValue::Boolean(v),
            _ => panic!("unsupported data type"),
        }
    }
}

impl From<&(opentelemetry::Key, opentelemetry::logs::AnyValue)> for LogRecordAttribute {
    fn from(value: &(opentelemetry::Key, opentelemetry::logs::AnyValue)) -> Self {
        Self {
            key: value.0.to_string(),
            value: value.1.to_owned().into(),
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

fn map_to_kv_list(
    map: Box<HashMap<opentelemetry::Key, opentelemetry::logs::AnyValue>>,
) -> Vec<LogMapKeyValue> {
    use opentelemetry::logs::AnyValue as A;
    map.into_iter()
        .map(|(k, val)| {
            let value = match val {
                A::Int(v) => LogMapValue::Value(v.into()),
                A::Double(v) => LogMapValue::Value(v.into()),
                A::String(v) => LogMapValue::Value(v.into()),
                A::Boolean(v) => LogMapValue::Value(v.into()),
                A::Bytes(l) => LogMapValue::List(l.into()),
                A::ListAny(l) => LogMapValue::List(l.into()),
                _ => panic!("unsupported data type"),
            };

            return LogMapKeyValue {
                key: k.to_string(),
                value,
            };
        })
        .collect()
}

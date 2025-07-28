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

impl From<opentelemetry::KeyValue> for KeyValue {
    fn from(value: opentelemetry::KeyValue) -> Self {
        Self {
            key: value.key.into(),
            value: value.value.into(),
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
            name: value.name.to_string(),
            description: value.description.to_string(),
            unit: value.unit.to_string(),
            data: (&value.data).into(),
        }
    }
}

impl From<Box<dyn opentelemetry_sdk::metrics::data::Aggregation>> for MetricData {
    fn from(value: Box<dyn opentelemetry_sdk::metrics::data::Aggregation>) -> Self {
        let v = value.as_any();
        if let Some(g) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<f64>>() {
            MetricData::Gauge(g.into())
        } else if let Some(g) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<i64>>() {
            MetricData::Gauge(g.into())
        } else if let Some(g) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<u64>>() {
            MetricData::Gauge(g.into())
        } else if let Some(s) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<f64>>() {
            MetricData::Sum(s.into())
        } else if let Some(s) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<i64>>() {
            MetricData::Sum(s.into())
        } else if let Some(s) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<u64>>() {
            MetricData::Sum(s.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Histogram<f64>>() {
            MetricData::Histogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Histogram<i64>>() {
            MetricData::Histogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Histogram<u64>>() {
            MetricData::Histogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::ExponentialHistogram<f64>>() {
            MetricData::ExponentialHistogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::ExponentialHistogram<i64>>() {
            MetricData::ExponentialHistogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::ExponentialHistogram<u64>>() {
            MetricData::ExponentialHistogram(h.into())
        } else {
            panic!("unknown aggregation type")
        }
    }
}

impl From<&Box<dyn opentelemetry_sdk::metrics::data::Aggregation>> for MetricData {
    fn from(value: &Box<dyn opentelemetry_sdk::metrics::data::Aggregation>) -> Self {
        let v = value.as_any();
        if let Some(g) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<f64>>() {
            MetricData::Gauge(g.into())
        } else if let Some(g) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<i64>>() {
            MetricData::Gauge(g.into())
        } else if let Some(g) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<u64>>() {
            MetricData::Gauge(g.into())
        } else if let Some(s) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<f64>>() {
            MetricData::Sum(s.into())
        } else if let Some(s) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<i64>>() {
            MetricData::Sum(s.into())
        } else if let Some(s) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<u64>>() {
            MetricData::Sum(s.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Histogram<f64>>() {
            MetricData::Histogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Histogram<i64>>() {
            MetricData::Histogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::Histogram<u64>>() {
            MetricData::Histogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::ExponentialHistogram<f64>>() {
            MetricData::ExponentialHistogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::ExponentialHistogram<i64>>() {
            MetricData::ExponentialHistogram(h.into())
        } else if let Some(h) = v.downcast_ref::<opentelemetry_sdk::metrics::data::ExponentialHistogram<u64>>() {
            MetricData::ExponentialHistogram(h.into())
        } else {
            panic!("unknown aggregation type")
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::Gauge<T>> for Gauge
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::Gauge<T>) -> Self {
        Self {
            data_points: value.data_points.to_owned().into_iter().map(Into::into).collect(),
            start_time: match value.start_time {
                Some(v) => Some(v.into()),
                None => None,
            },
            time: value.time.into(),
        }
    }
}

impl<T> From<opentelemetry_sdk::metrics::data::GaugeDataPoint<T>> for GaugeDataPoint
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: opentelemetry_sdk::metrics::data::GaugeDataPoint<T>) -> Self {
        Self {
            attributes: value.attributes.into_iter().map(Into::into).collect(),
            value: value.value.into(),
            exemplars: value.exemplars.into_iter().map(Into::into).collect(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::Sum<T>> for Sum
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::Sum<T>) -> Self {
        Self {
            data_points: value.data_points.to_owned().into_iter().map(Into::into).collect(),
            start_time: value.start_time.into(),
            time: value.time.into(),
            temporality: value.temporality.into(),
            is_monotonic: value.is_monotonic,
        }
    }
}

impl<T> From<opentelemetry_sdk::metrics::data::SumDataPoint<T>> for SumDataPoint
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: opentelemetry_sdk::metrics::data::SumDataPoint<T>) -> Self {
        Self {
            attributes: value.attributes.into_iter().map(Into::into).collect(),
            value: value.value.into(),
            exemplars: value.exemplars.into_iter().map(Into::into).collect(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::Histogram<T>> for Histogram
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::Histogram<T>) -> Self {
        Self {
            data_points: value.data_points.to_owned().into_iter().map(Into::into).collect(),
            start_time: value.start_time.into(),
            time: value.time.into(),
            temporality: value.temporality.into(),
        }
    }
}

impl<T> From<opentelemetry_sdk::metrics::data::HistogramDataPoint<T>> for HistogramDataPoint
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: opentelemetry_sdk::metrics::data::HistogramDataPoint<T>) -> Self {
        Self {
            attributes: value.attributes.into_iter().map(Into::into).collect(),
            count: value.count,
            bounds: value.bounds,
            bucket_counts: value.bucket_counts,
            min: match value.min {
                Some(v) => Some(v.into()),
                None => None,
            },
            max: match value.max {
                Some(v) => Some(v.into()),
                None => None,
            },
            sum: value.sum.into(),
            exemplars: value.exemplars.into_iter().map(Into::into).collect(),
        }
    }
}


impl<T> From<&opentelemetry_sdk::metrics::data::ExponentialHistogram<T>> for ExponentialHistogram
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::ExponentialHistogram<T>) -> Self {
        Self {
            data_points: value.data_points.to_owned().into_iter().map(Into::into).collect(),
            start_time: value.start_time.into(),
            time: value.time.into(),
            temporality: value.temporality.into(),
        }
    }
}

impl<T> From<opentelemetry_sdk::metrics::data::ExponentialHistogramDataPoint<T>> for ExponentialHistogramDataPoint
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: opentelemetry_sdk::metrics::data::ExponentialHistogramDataPoint<T>) -> Self {
        Self {
            attributes: value.attributes.into_iter().map(Into::into).collect(),
            count: value.count as u64,
            min: match value.min {
                Some(v) => Some(v.into()),
                None => None,
            },
            max: match value.max {
                Some(v) => Some(v.into()),
                None => None,
            },
            sum: value.sum.into(),
            scale: value.scale,
            zero_count: value.zero_count,
            positive_bucket: value.positive_bucket.into(),
            negative_bucket: value.negative_bucket.into(),
            zero_threshold: value.zero_threshold,
            exemplars: value.exemplars.into_iter().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry_sdk::metrics::data::ExponentialBucket> for ExponentialBucket {
    fn from(value: opentelemetry_sdk::metrics::data::ExponentialBucket) -> Self {
        Self {
            offset: value.offset,
            counts: value.counts,
        }
    }
}

impl From<opentelemetry_sdk::metrics::Temporality> for TemporalityT {
    fn from(value: opentelemetry_sdk::metrics::Temporality) -> Self {
        match value {
            opentelemetry_sdk::metrics::Temporality::Delta => TemporalityT::Delta,
            opentelemetry_sdk::metrics::Temporality::LowMemory => TemporalityT::LowMemory,
            _ => TemporalityT::Cumulative,
        }
    }
}

impl<T> From<opentelemetry_sdk::metrics::data::Exemplar<T>> for Exemplar
where
    T: Into<MetricNumber>
{
     fn from(value: opentelemetry_sdk::metrics::data::Exemplar<T>) -> Self {
         Self {
            filtered_attributes: value.filtered_attributes.into_iter().map(Into::into).collect(),
            time: value.time.into(),
            value: value.value.into(),
            span_id: String::from_utf8(value.span_id.to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id.to_vec()).unwrap(),
         }
     }
}

impl From<f64> for MetricNumber {
    fn from(value: f64) -> Self {
        MetricNumber::F64(value)
    }
}

impl From<i64> for MetricNumber {
    fn from(value: i64) -> Self {
        MetricNumber::S64(value)
    }
}

impl From<u64> for MetricNumber {
    fn from(value: u64) -> Self {
        MetricNumber::U64(value)
    }
}
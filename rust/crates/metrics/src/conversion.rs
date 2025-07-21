use crate::wit::wasi::otel::metrics::*;

impl From<opentelemetry_sdk::metrics::data::ResourceMetrics> for ResourceMetrics {
    fn from(value: opentelemetry_sdk::metrics::data::ResourceMetrics) -> Self {
        Self {
            resource: value.resource.into(),
            scope_metrics: value.scope_metrics.into_iter().map(Into::into).collect(),
        }
    }
}

impl From<&mut opentelemetry_sdk::metrics::data::ResourceMetrics> for ResourceMetrics {
    fn from(value: &mut opentelemetry_sdk::metrics::data::ResourceMetrics) -> Self {
        Self {
            resource: (&value.resource).into(),
            scope_metrics: value.scope_metrics.iter().to_owned().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::ResourceMetrics> for ResourceMetrics {
    fn from(value: &opentelemetry_sdk::metrics::data::ResourceMetrics) -> Self {
        Self {
            resource: (&value.resource).into(),
            scope_metrics: value.scope_metrics.iter().to_owned().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry_sdk::resource::Resource> for Resource {
    fn from(value: opentelemetry_sdk::resource::Resource) -> Self {
        let mut attrs: Vec<KeyValue> = Vec::new();
        value.into_iter().for_each(|v| {
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

impl From<&opentelemetry_sdk::resource::Resource> for Resource {
    fn from(value: &opentelemetry_sdk::resource::Resource) -> Self {
        let mut attrs: Vec<KeyValue> = Vec::new();
        value.into_iter().for_each(|v| {
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

impl From<&opentelemetry_sdk::metrics::data::ScopeMetrics> for ScopeMetrics {
    fn from(value: &opentelemetry_sdk::metrics::data::ScopeMetrics) -> Self {
        Self {
            scope: value.scope.clone().into(),
            metrics: value.metrics.iter().to_owned().map(Into::into).collect(),
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

impl From<&opentelemetry_sdk::metrics::data::Metric> for Metric {
    fn from(value: &opentelemetry_sdk::metrics::data::Metric) -> Self {
        Self {
            name: value.name.to_owned().into(),
            description: value.description.to_owned().into(),
            unit: value.unit.to_owned().into(),
            data: convert_aggregation(&value.data),
        }
    }
}

fn convert_aggregation(v: &Box<dyn opentelemetry_sdk::metrics::data::Aggregation>) -> Aggregation {
    let any_ref = v.as_any();
    if let Some(gauge_f64) = any_ref.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<f64>>() {
        Aggregation::Gauge(Gauge { data_points: gauge_f64.data_points.iter().map(Into::into).collect() })
    } else if let Some(gauge_i64) = any_ref.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<i64>>() {
        Aggregation::Gauge(Gauge { data_points: gauge_i64.data_points.iter().map(Into::into).collect() })
    } else if let Some(gauge_u64) = any_ref.downcast_ref::<opentelemetry_sdk::metrics::data::Gauge<u64>>() {
        Aggregation::Gauge(Gauge { data_points: gauge_u64.data_points.iter().map(Into::into).collect() })
    } else if let Some(sum_f64) = any_ref.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<f64>>() {
        Aggregation::Sum(Sum {
            data_points: sum_f64.data_points.iter().map(Into::into).collect(),
            temporality: sum_f64.temporality.into(),
            is_monotonic: sum_f64.is_monotonic,
        })
    } else if let Some(sum_i64) = any_ref.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<i64>>() {
        Aggregation::Sum(Sum {
            data_points: sum_i64.data_points.iter().map(Into::into).collect(),
            temporality: sum_i64.temporality.into(),
            is_monotonic: sum_i64.is_monotonic,
        })
    } else if let Some(sum_u64) = any_ref.downcast_ref::<opentelemetry_sdk::metrics::data::Sum<u64>>() {
        Aggregation::Sum(Sum {
            data_points: sum_u64.data_points.iter().map(Into::into).collect(),
            temporality: sum_u64.temporality.into(),
            is_monotonic: sum_u64.is_monotonic,
        })
    } else if let Some(histogram_f64) = any_ref.downcast_ref::<opentelemetry_sdk::metrics::data::Histogram<f64>>() {
        Aggregation::Histogram(Histogram {
            data_points: histogram_f64.data_points.iter().map(Into::into).collect(),
            temporality: histogram_f64.temporality.into(),
        })
    } else if let Some(histogram_u64) = any_ref.downcast_ref::<opentelemetry_sdk::metrics::data::Histogram<u64>>() {
        Aggregation::Histogram(Histogram {
            data_points: histogram_u64.data_points.iter().map(Into::into).collect(),
            temporality: histogram_u64.temporality.into(),
        })
    } else {
        // TODO: handle this better
        panic!("ERROR")
    }
}

impl From<opentelemetry_sdk::metrics::Temporality> for TemporalityT {
    fn from(value: opentelemetry_sdk::metrics::Temporality) -> Self {
        match value {
            opentelemetry_sdk::metrics::Temporality::Cumulative => TemporalityT::Cumulative,
            opentelemetry_sdk::metrics::Temporality::Delta => TemporalityT::Delta,
            _ => TemporalityT::LowMemory,
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::DataPoint<f64>> for DataPoint {
    fn from(value: &opentelemetry_sdk::metrics::data::DataPoint<f64>) -> Self {
        Self {
            attributes: value.attributes.to_owned().into_iter().map(Into::into).collect(),
            start_time: match value.start_time {
                Some(v) => Some(v.into()),
                None => None,
            },
            time: match value.time {
                Some(v) => Some(v.into()),
                None => None,
            },
            value: AggregationNumber::F64(value.value),
            exemplars: value.exemplars.to_owned().into_iter().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::DataPoint<i64>> for DataPoint {
    fn from(value: &opentelemetry_sdk::metrics::data::DataPoint<i64>) -> Self {
        Self {
            attributes: value.attributes.to_owned().into_iter().map(Into::into).collect(),
            start_time: match value.start_time {
                Some(v) => Some(v.into()),
                None => None,
            },
            time: match value.time {
                Some(v) => Some(v.into()),
                None => None,
            },
            value: AggregationNumber::S64(value.value),
            exemplars: value.exemplars.to_owned().into_iter().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::DataPoint<u64>> for DataPoint {
    fn from(value: &opentelemetry_sdk::metrics::data::DataPoint<u64>) -> Self {
        Self {
            attributes: value.attributes.to_owned().into_iter().map(Into::into).collect(),
            start_time: match value.start_time {
                Some(v) => Some(v.into()),
                None => None,
            },
            time: match value.time {
                Some(v) => Some(v.into()),
                None => None,
            },
            value: AggregationNumber::U64(value.value),
            exemplars: value.exemplars.to_owned().into_iter().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::HistogramDataPoint<f64>> for HistogramDataPoint {
    fn from(value: &opentelemetry_sdk::metrics::data::HistogramDataPoint<f64>) -> Self {
        Self {
            attributes: value.attributes.iter().map(Into::into).collect(),
            start_time: value.start_time.into(),
            time: value.time.into(),
            count: value.count,
            bounds: value.bounds.to_owned(),
            bucket_counts: value.bucket_counts.to_owned(),
            min: match value.min {
                Some(v) => Some(AggregationNumber::F64(v)),
                None => None,
            },
            max: match value.max {
                Some(v) => Some(AggregationNumber::F64(v)),
                None => None,
            },
            sum: AggregationNumber::F64(value.sum),
            exemplars: value.exemplars.iter().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::HistogramDataPoint<i64>> for HistogramDataPoint {
    fn from(value: &opentelemetry_sdk::metrics::data::HistogramDataPoint<i64>) -> Self {
        Self {
            attributes: value.attributes.iter().map(Into::into).collect(),
            start_time: value.start_time.into(),
            time: value.time.into(),
            count: value.count,
            bounds: value.bounds.to_owned(),
            bucket_counts: value.bucket_counts.to_owned(),
            min: match value.min {
                Some(v) => Some(AggregationNumber::S64(v)),
                None => None,
            },
            max: match value.max {
                Some(v) => Some(AggregationNumber::S64(v)),
                None => None,
            },
            sum: AggregationNumber::S64(value.sum),
            exemplars: value.exemplars.iter().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::HistogramDataPoint<u64>> for HistogramDataPoint {
    fn from(value: &opentelemetry_sdk::metrics::data::HistogramDataPoint<u64>) -> Self {
        Self {
            attributes: value.attributes.iter().map(Into::into).collect(),
            start_time: value.start_time.into(),
            time: value.time.into(),
            count: value.count,
            bounds: value.bounds.to_owned(),
            bucket_counts: value.bucket_counts.to_owned(),
            min: match value.min {
                Some(v) => Some(AggregationNumber::U64(v)),
                None => None,
            },
            max: match value.max {
                Some(v) => Some(AggregationNumber::U64(v)),
                None => None,
            },
            sum: AggregationNumber::U64(value.sum),
            exemplars: value.exemplars.iter().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Exemplar<f64>> for Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<f64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes.iter().map(Into::into).collect(),
            time: value.time.into(),
            value: AggregationNumber::F64(value.value),
            span_id: String::from_utf8(value.span_id.to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id.to_vec()).unwrap(),
        }
    }
}

impl From<opentelemetry_sdk::metrics::data::Exemplar<f64>> for Exemplar {
    fn from(value: opentelemetry_sdk::metrics::data::Exemplar<f64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes.iter().map(Into::into).collect(),
            time: value.time.into(),
            value: AggregationNumber::F64(value.value),
            span_id: String::from_utf8(value.span_id.to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id.to_vec()).unwrap(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Exemplar<i64>> for Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<i64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes.iter().map(Into::into).collect(),
            time: value.time.into(),
            value: AggregationNumber::S64(value.value),
            span_id: String::from_utf8(value.span_id.to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id.to_vec()).unwrap(),
        }
    }
}

impl From<opentelemetry_sdk::metrics::data::Exemplar<i64>> for Exemplar {
    fn from(value: opentelemetry_sdk::metrics::data::Exemplar<i64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes.iter().map(Into::into).collect(),
            time: value.time.into(),
            value: AggregationNumber::S64(value.value),
            span_id: String::from_utf8(value.span_id.to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id.to_vec()).unwrap(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Exemplar<u64>> for Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<u64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes.iter().map(Into::into).collect(),
            time: value.time.into(),
            value: AggregationNumber::U64(value.value),
            span_id: String::from_utf8(value.span_id.to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id.to_vec()).unwrap(),
        }
    }
}

impl From<opentelemetry_sdk::metrics::data::Exemplar<u64>> for Exemplar {
    fn from(value: opentelemetry_sdk::metrics::data::Exemplar<u64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes.iter().map(Into::into).collect(),
            time: value.time.into(),
            value: AggregationNumber::U64(value.value),
            span_id: String::from_utf8(value.span_id.to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id.to_vec()).unwrap(),
        }
    }
}

impl From<Box<dyn opentelemetry_sdk::metrics::data::Aggregation>> for Aggregation {
    fn from(value: Box<dyn opentelemetry_sdk::metrics::data::Aggregation>) -> Self {
        value.into()
    }
}

impl From<Box<&dyn opentelemetry_sdk::metrics::data::Aggregation>> for Aggregation {
    fn from(value: Box<&dyn opentelemetry_sdk::metrics::data::Aggregation>) -> Self {
        value.into()
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

// TODO: duplicated code
impl From<std::time::SystemTime> for Datetime {
    fn from(value: std::time::SystemTime) -> Self {
        let duration_since_epoch = value
            .duration_since(std::time::UNIX_EPOCH)
            .expect("SystemTime should be after UNIX EPOCH");
        Self {
            seconds: duration_since_epoch.as_secs(),
            nanoseconds: duration_since_epoch.subsec_nanos(),
        }
    }
}

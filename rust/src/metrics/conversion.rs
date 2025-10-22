use crate::wit::wasi::otel::metrics::*;

impl From<opentelemetry_sdk::metrics::data::ResourceMetrics> for ResourceMetrics {
    fn from(value: opentelemetry_sdk::metrics::data::ResourceMetrics) -> Self {
        Self {
            resource: value.resource().to_owned().into(),
            scope_metrics: value.scope_metrics().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::ScopeMetrics> for ScopeMetrics {
    fn from(value: &opentelemetry_sdk::metrics::data::ScopeMetrics) -> Self {
        Self {
            scope: value.scope().into(),
            metrics: value.metrics().map(Into::into).collect(),
        }
    }
}

impl From<opentelemetry_sdk::resource::Resource> for Resource {
    fn from(value: opentelemetry_sdk::resource::Resource) -> Self {
        Self {
            attributes: value.into_iter().map(Into::into).collect(),
            schema_url: match value.schema_url() {
                Some(v) => Some(v.to_string()),
                None => None,
            },
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Metric> for Metric {
    fn from(value: &opentelemetry_sdk::metrics::data::Metric) -> Self {
        Self {
            name: value.name().to_string(),
            description: value.description().to_string(),
            unit: value.unit().to_string(),
            data: value.data().into(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::AggregatedMetrics> for AggregatedMetrics {
    fn from(value: &opentelemetry_sdk::metrics::data::AggregatedMetrics) -> Self {
        use opentelemetry_sdk::metrics::data as sdk;
        fn convert_metric<T>(metric: &opentelemetry_sdk::metrics::data::MetricData<T>) -> MetricData
        where
            T: Into<MetricNumber> + Clone + Copy,
        {
            match metric {
                sdk::MetricData::Gauge(gauge) => MetricData::Gauge(Gauge {
                    data_points: gauge.data_points().map(Into::into).collect(),
                    start_time: match gauge.start_time() {
                        Some(v) => Some(v.into()),
                        None => None,
                    },
                    time: gauge.time().into(),
                }),
                sdk::MetricData::Sum(sum) => MetricData::Sum(Sum {
                    data_points: sum.data_points().map(Into::into).collect(),
                    start_time: sum.start_time().into(),
                    time: sum.time().into(),
                    temporality: sum.temporality().into(),
                    is_monotonic: sum.is_monotonic(),
                }),
                sdk::MetricData::Histogram(hist) => MetricData::Histogram(Histogram {
                    data_points: hist.data_points().map(Into::into).collect(),
                    start_time: hist.start_time().into(),
                    time: hist.time().into(),
                    temporality: hist.temporality().into(),
                }),
                sdk::MetricData::ExponentialHistogram(hist) => {
                    MetricData::ExponentialHistogram(ExponentialHistogram {
                        data_points: hist.data_points().map(Into::into).collect(),
                        start_time: hist.start_time().into(),
                        time: hist.time().into(),
                        temporality: hist.temporality().into(),
                    })
                }
            }
        }
        match value {
            sdk::AggregatedMetrics::F64(v) => AggregatedMetrics::F64(convert_metric(v)),
            sdk::AggregatedMetrics::U64(v) => AggregatedMetrics::U64(convert_metric(v)),
            sdk::AggregatedMetrics::I64(v) => AggregatedMetrics::S64(convert_metric(v)),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::Gauge<T>> for Gauge
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::Gauge<T>) -> Self {
        Self {
            data_points: value.data_points().map(Into::into).collect(),
            start_time: match value.start_time() {
                Some(v) => Some(v.into()),
                None => None,
            },
            time: value.time().into(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::GaugeDataPoint<T>> for GaugeDataPoint
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::GaugeDataPoint<T>) -> Self {
        Self {
            attributes: value.attributes().into_iter().map(Into::into).collect(),
            value: value.value().into(),
            exemplars: value.exemplars().into_iter().map(Into::into).collect(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::Sum<T>> for Sum
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::Sum<T>) -> Self {
        Self {
            data_points: value.data_points().map(Into::into).collect(),
            start_time: value.start_time().into(),
            time: value.time().into(),
            temporality: value.temporality().into(),
            is_monotonic: value.is_monotonic(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::SumDataPoint<T>> for SumDataPoint
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::SumDataPoint<T>) -> Self {
        Self {
            attributes: value.attributes().into_iter().map(Into::into).collect(),
            value: value.value().into(),
            exemplars: value.exemplars().into_iter().map(Into::into).collect(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::Histogram<T>> for Histogram
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::Histogram<T>) -> Self {
        Self {
            data_points: value.data_points().map(Into::into).collect(),
            start_time: value.start_time().into(),
            time: value.time().into(),
            temporality: value.temporality().into(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::HistogramDataPoint<T>> for HistogramDataPoint
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::HistogramDataPoint<T>) -> Self {
        Self {
            attributes: value.attributes().into_iter().map(Into::into).collect(),
            count: value.count(),
            bounds: value.bounds().collect(),
            bucket_counts: value.bucket_counts().collect(),
            min: match value.min() {
                Some(v) => Some(v.into()),
                None => None,
            },
            max: match value.max() {
                Some(v) => Some(v.into()),
                None => None,
            },
            sum: value.sum().into(),
            exemplars: value.exemplars().into_iter().map(Into::into).collect(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::ExponentialHistogram<T>> for ExponentialHistogram
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::ExponentialHistogram<T>) -> Self {
        Self {
            data_points: value.data_points().map(Into::into).collect(),
            start_time: value.start_time().into(),
            time: value.time().into(),
            temporality: value.temporality().into(),
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::ExponentialHistogramDataPoint<T>>
    for ExponentialHistogramDataPoint
where
    T: Into<MetricNumber> + Clone + Copy,
{
    fn from(value: &opentelemetry_sdk::metrics::data::ExponentialHistogramDataPoint<T>) -> Self {
        Self {
            attributes: value.attributes().into_iter().map(Into::into).collect(),
            count: value.count() as u64,
            min: match value.min() {
                Some(v) => Some(v.into()),
                None => None,
            },
            max: match value.max() {
                Some(v) => Some(v.into()),
                None => None,
            },
            sum: value.sum().into(),
            scale: value.scale(),
            zero_count: value.zero_count(),
            positive_bucket: value.positive_bucket().into(),
            negative_bucket: value.negative_bucket().into(),
            zero_threshold: value.zero_threshold(),
            exemplars: value.exemplars().map(Into::into).collect(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::ExponentialBucket> for ExponentialBucket {
    fn from(value: &opentelemetry_sdk::metrics::data::ExponentialBucket) -> Self {
        Self {
            offset: value.offset(),
            counts: value.counts().collect(),
        }
    }
}

impl From<opentelemetry_sdk::metrics::Temporality> for Temporality {
    fn from(value: opentelemetry_sdk::metrics::Temporality) -> Self {
        match value {
            opentelemetry_sdk::metrics::Temporality::Delta => Temporality::Delta,
            opentelemetry_sdk::metrics::Temporality::LowMemory => Temporality::LowMemory,
            _ => Temporality::Cumulative,
        }
    }
}

impl<T> From<&opentelemetry_sdk::metrics::data::Exemplar<T>> for Exemplar
where
    T: Into<MetricNumber> + Clone,
{
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<T>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes().map(Into::into).collect(),
            time: value.time().into(),
            value: value.clone().value.into(), // TODO: this feels heavy...research optimizing?
            span_id: String::from_utf8(value.span_id().to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id().to_vec()).unwrap(),
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

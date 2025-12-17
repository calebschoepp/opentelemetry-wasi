use crate::wit::wasi::otel::metrics::*;

impl From<opentelemetry_sdk::metrics::data::ResourceMetrics> for ResourceMetrics {
    fn from(value: opentelemetry_sdk::metrics::data::ResourceMetrics) -> Self {
        Self {
            resource: value.resource().into(),
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

// Convert OTel MetricData to WASI MetricData
macro_rules! metric_data_to_wasi {
    (
        $otel_metric_data:expr,
        $wasi_gauge_type:ident,
        $wasi_sum_type:ident,
        $wasi_histogram_type:ident,
        $wasi_exponential_histogram_type:ident,
    ) => {
        match $otel_metric_data {
            opentelemetry_sdk::metrics::data::MetricData::Gauge(g) => {
                crate::wit::wasi::otel::metrics::MetricData::$wasi_gauge_type(
                    crate::wit::wasi::otel::metrics::Gauge {
                        data_points: g
                            .data_points()
                            .map(|dp| crate::wit::wasi::otel::metrics::GaugeDataPoint {
                                attributes: dp.attributes().into_iter().map(Into::into).collect(),
                                value: dp.value().into(),
                                exemplars: dp.exemplars().into_iter().map(Into::into).collect(),
                            })
                            .collect(),
                        start_time: g.start_time().map(Into::into),
                        time: g.time().into(),
                    },
                )
            }
            opentelemetry_sdk::metrics::data::MetricData::Sum(s) => {
                crate::wit::wasi::otel::metrics::MetricData::$wasi_sum_type(
                    crate::wit::wasi::otel::metrics::Sum {
                        data_points: s
                            .data_points()
                            .map(|dp| crate::wit::wasi::otel::metrics::SumDataPoint {
                                attributes: dp.attributes().into_iter().map(Into::into).collect(),
                                value: dp.value().into(),
                                exemplars: dp.exemplars().into_iter().map(Into::into).collect(),
                            })
                            .collect(),
                        start_time: s.start_time().into(),
                        time: s.time().into(),
                        temporality: s.temporality().into(),
                        is_monotonic: s.is_monotonic(),
                    },
                )
            }
            opentelemetry_sdk::metrics::data::MetricData::Histogram(h) => {
                crate::wit::wasi::otel::metrics::MetricData::$wasi_histogram_type(
                    crate::wit::wasi::otel::metrics::Histogram {
                        data_points: h
                            .data_points()
                            .map(|dp| crate::wit::wasi::otel::metrics::HistogramDataPoint {
                                attributes: dp.attributes().into_iter().map(Into::into).collect(),
                                count: dp.count(),
                                bounds: dp.bounds().collect(),
                                bucket_counts: dp.bucket_counts().collect(),
                                min: dp.min().map(Into::into),
                                max: dp.max().map(Into::into),
                                sum: dp.sum().into(),
                                exemplars: dp.exemplars().into_iter().map(Into::into).collect(),
                            })
                            .collect(),
                        start_time: h.start_time().into(),
                        time: h.time().into(),
                        temporality: h.temporality().into(),
                    },
                )
            }
            opentelemetry_sdk::metrics::data::MetricData::ExponentialHistogram(h) => {
                crate::wit::wasi::otel::metrics::MetricData::$wasi_exponential_histogram_type(
                    crate::wit::wasi::otel::metrics::ExponentialHistogram {
                        data_points: h
                            .data_points()
                            .map(|dp| {
                                crate::wit::wasi::otel::metrics::ExponentialHistogramDataPoint {
                                    attributes: dp
                                        .attributes()
                                        .into_iter()
                                        .map(Into::into)
                                        .collect(),
                                    count: dp.count() as u64,
                                    min: dp.min().map(Into::into),
                                    max: dp.max().map(Into::into),
                                    sum: dp.sum().into(),
                                    scale: dp.scale(),
                                    zero_count: dp.zero_count(),
                                    positive_bucket: dp.positive_bucket().into(),
                                    negative_bucket: dp.negative_bucket().into(),
                                    zero_threshold: dp.zero_threshold(),
                                    exemplars: dp.exemplars().map(Into::into).collect(),
                                }
                            })
                            .collect(),
                        start_time: h.start_time().into(),
                        time: h.time().into(),
                        temporality: h.temporality().into(),
                    },
                )
            }
        }
    };
}

impl From<&opentelemetry_sdk::metrics::data::AggregatedMetrics> for MetricData {
    fn from(value: &opentelemetry_sdk::metrics::data::AggregatedMetrics) -> Self {
        use opentelemetry_sdk::metrics::data as sdk;
        match value {
            sdk::AggregatedMetrics::F64(data) => metric_data_to_wasi!(
                data,
                F64Gauge,
                F64Sum,
                F64Histogram,
                F64ExponentialHistogram,
            ),
            sdk::AggregatedMetrics::U64(data) => metric_data_to_wasi!(
                data,
                U64Gauge,
                U64Sum,
                U64Histogram,
                U64ExponentialHistogram,
            ),
            sdk::AggregatedMetrics::I64(data) => metric_data_to_wasi!(
                data,
                S64Gauge,
                S64Sum,
                S64Histogram,
                S64ExponentialHistogram,
            ),
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

impl From<&opentelemetry_sdk::metrics::data::Exemplar<f64>> for Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<f64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes().map(Into::into).collect(),
            time: value.time().into(),
            value: value.value.into(),
            span_id: String::from_utf8(value.span_id().to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id().to_vec()).unwrap(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Exemplar<u64>> for Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<u64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes().map(Into::into).collect(),
            time: value.time().into(),
            value: value.value.into(),
            span_id: String::from_utf8(value.span_id().to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id().to_vec()).unwrap(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Exemplar<i64>> for Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<i64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes().map(Into::into).collect(),
            time: value.time().into(),
            value: value.value.into(),
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

impl From<u64> for MetricNumber {
    fn from(value: u64) -> Self {
        MetricNumber::U64(value)
    }
}

impl From<i64> for MetricNumber {
    fn from(value: i64) -> Self {
        MetricNumber::S64(value)
    }
}

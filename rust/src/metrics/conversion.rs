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

impl From<&opentelemetry_sdk::metrics::data::AggregatedMetrics> for MetricData {
    fn from(value: &opentelemetry_sdk::metrics::data::AggregatedMetrics) -> Self {
        use opentelemetry_sdk::metrics::data as sdk;
        match value {
            sdk::AggregatedMetrics::F64(data) => match data {
                sdk::MetricData::Gauge(g) => MetricData::F64Gauge(F64Gauge {
                    data_points: g
                        .data_points()
                        .map(|dp| F64GaugeDataPoint {
                            attributes: dp.attributes().into_iter().map(Into::into).collect(),
                            value: dp.value().into(),
                            exemplars: dp.exemplars().into_iter().map(Into::into).collect(),
                        })
                        .collect(),
                    start_time: match g.start_time() {
                        Some(v) => Some(v.into()),
                        None => None,
                    },
                    time: g.time().into(),
                }),
                sdk::MetricData::Sum(s) => MetricData::F64Sum(F64Sum {
                    data_points: s
                        .data_points()
                        .map(|dp| F64SumDataPoint {
                            attributes: dp.attributes().into_iter().map(Into::into).collect(),
                            value: dp.value().into(),
                            exemplars: dp.exemplars().into_iter().map(Into::into).collect(),
                        })
                        .collect(),
                    start_time: s.start_time().into(),
                    time: s.time().into(),
                    temporality: s.temporality().into(),
                    is_monotonic: s.is_monotonic(),
                }),
                sdk::MetricData::Histogram(h) => MetricData::F64Histogram(F64Histogram {
                    data_points: h
                        .data_points()
                        .map(|dp| F64HistogramDataPoint {
                            attributes: dp.attributes().into_iter().map(Into::into).collect(),
                            count: dp.count(),
                            bounds: dp.bounds().collect(),
                            bucket_counts: dp.bucket_counts().collect(),
                            min: match dp.min() {
                                Some(v) => Some(v.into()),
                                None => None,
                            },
                            max: match dp.max() {
                                Some(v) => Some(v.into()),
                                None => None,
                            },
                            sum: dp.sum().into(),
                            exemplars: dp.exemplars().into_iter().map(Into::into).collect(),
                        })
                        .collect(),
                    start_time: h.start_time().into(),
                    time: h.time().into(),
                    temporality: h.temporality().into(),
                }),
                sdk::MetricData::ExponentialHistogram(h) => {
                    MetricData::F64ExponentialHistogram(F64ExponentialHistogram {
                        data_points: h
                            .data_points()
                            .map(|dp| F64ExponentialHistogramDataPoint {
                                attributes: dp.attributes().into_iter().map(Into::into).collect(),
                                count: dp.count() as u64,
                                min: match dp.min() {
                                    Some(v) => Some(v.into()),
                                    None => None,
                                },
                                max: match dp.max() {
                                    Some(v) => Some(v.into()),
                                    None => None,
                                },
                                sum: dp.sum().into(),
                                scale: dp.scale(),
                                zero_count: dp.zero_count(),
                                positive_bucket: dp.positive_bucket().into(),
                                negative_bucket: dp.negative_bucket().into(),
                                zero_threshold: dp.zero_threshold(),
                                exemplars: dp.exemplars().map(Into::into).collect(),
                            })
                            .collect(),
                        start_time: h.start_time().into(),
                        time: h.time().into(),
                        temporality: h.temporality().into(),
                    })
                }
            },
            _ => todo!("Create a macro for the remaining number types"),
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

impl From<&opentelemetry_sdk::metrics::data::Exemplar<f64>> for F64Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<f64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes().map(Into::into).collect(),
            time: value.time().into(),
            value: value.value,
            span_id: String::from_utf8(value.span_id().to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id().to_vec()).unwrap(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Exemplar<u64>> for U64Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<u64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes().map(Into::into).collect(),
            time: value.time().into(),
            value: value.value,
            span_id: String::from_utf8(value.span_id().to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id().to_vec()).unwrap(),
        }
    }
}

impl From<&opentelemetry_sdk::metrics::data::Exemplar<i64>> for S64Exemplar {
    fn from(value: &opentelemetry_sdk::metrics::data::Exemplar<i64>) -> Self {
        Self {
            filtered_attributes: value.filtered_attributes().map(Into::into).collect(),
            time: value.time().into(),
            value: value.value,
            span_id: String::from_utf8(value.span_id().to_vec()).unwrap(),
            trace_id: String::from_utf8(value.trace_id().to_vec()).unwrap(),
        }
    }
}

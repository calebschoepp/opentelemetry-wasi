package metrics

import (
	"time"

	"github.com/calebschoepp/opentelemetry-wasi/types"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_clocks_wall_clock"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_metrics"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_types"
	"github.com/calebschoepp/opentelemetry-wasi/wit_component/wit_types"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func toWasiResourceMetrics(rm metricdata.ResourceMetrics) wasi_otel_metrics.ResourceMetrics {
	var r wasi_otel_types.Resource
	if rm.Resource != nil {
		r = types.ToWasiResource(*rm.Resource)
	}

	return wasi_otel_metrics.ResourceMetrics{
		Resource:     r, // Default to empty resource if nil
		ScopeMetrics: toWasiScopeMetrics(rm.ScopeMetrics),
	}
}

func toWasiScopeMetrics(sm []metricdata.ScopeMetrics) []wasi_otel_metrics.ScopeMetrics {
	result := make([]wasi_otel_metrics.ScopeMetrics, len(sm))

	for i, m := range sm {
		result[i] = wasi_otel_metrics.ScopeMetrics{
			Scope:   types.ToWasiInstrumentationScope(m.Scope),
			Metrics: toWasiMetrics(m.Metrics),
		}
	}

	return result
}

func toWasiMetrics(metrics []metricdata.Metrics) []wasi_otel_metrics.Metric {
	result := make([]wasi_otel_metrics.Metric, len(metrics))
	for i, m := range metrics {
		result[i] = wasi_otel_metrics.Metric{
			Name:        m.Name,
			Description: m.Description,
			Unit:        m.Unit,
			Data:        toWasiMetricData(m.Data),
		}
	}

	return result
}

func toWasiMetricData(data metricdata.Aggregation) wasi_otel_metrics.MetricData {
	switch v := data.(type) {
	case metricdata.Gauge[int64]:
		var startOpt wit_types.Option[wasi_clocks_wall_clock.Datetime]
		start, timeVal := extractTimestamps(v.DataPoints)
		if start == nil {
			startOpt = wit_types.None[wasi_clocks_wall_clock.Datetime]()
		} else {
			startOpt = wit_types.Some(*start)
		}

		return wasi_otel_metrics.MakeMetricDataS64Gauge(wasi_otel_metrics.Gauge{
			DataPoints: toWasiGaugeDataPoint(v.DataPoints),
			StartTime:  startOpt,
			Time:       timeVal,
		})
	case metricdata.Gauge[float64]:
		var startOpt wit_types.Option[wasi_clocks_wall_clock.Datetime]
		start, timeVal := extractTimestamps(v.DataPoints)
		if start == nil {
			startOpt = wit_types.None[wasi_clocks_wall_clock.Datetime]()
		} else {
			startOpt = wit_types.Some(*start)
		}
		return wasi_otel_metrics.MakeMetricDataF64Gauge(wasi_otel_metrics.Gauge{
			DataPoints: toWasiGaugeDataPoint(v.DataPoints),
			StartTime:  startOpt,
			Time:       timeVal,
		})
	case metricdata.Sum[int64]:
		start, timeVal := extractTimestamps(v.DataPoints)
		return wasi_otel_metrics.MakeMetricDataS64Sum(wasi_otel_metrics.Sum{
			DataPoints:  toWasiSumDataPoint(v.DataPoints),
			StartTime:   *start,
			Time:        timeVal,
			Temporality: toWasiTemporality(v.Temporality),
			IsMonotonic: v.IsMonotonic,
		})
	case metricdata.Sum[float64]:
		start, timeVal := extractTimestamps(v.DataPoints)
		return wasi_otel_metrics.MakeMetricDataF64Sum(wasi_otel_metrics.Sum{
			DataPoints:  toWasiSumDataPoint(v.DataPoints),
			StartTime:   *start,
			Time:        timeVal,
			Temporality: toWasiTemporality(v.Temporality),
			IsMonotonic: v.IsMonotonic,
		})
	case metricdata.ExponentialHistogram[int64]:
		start, timeVal := extractTimestamps(v.DataPoints)
		return wasi_otel_metrics.MakeMetricDataS64ExponentialHistogram(wasi_otel_metrics.ExponentialHistogram{
			DataPoints:  toWasiExponentialHistogram(v.DataPoints),
			StartTime:   *start,
			Time:        timeVal,
			Temporality: toWasiTemporality(v.Temporality),
		})
	case metricdata.ExponentialHistogram[float64]:
		start, timeVal := extractTimestamps(v.DataPoints)
		return wasi_otel_metrics.MakeMetricDataF64ExponentialHistogram(wasi_otel_metrics.ExponentialHistogram{
			DataPoints:  toWasiExponentialHistogram(v.DataPoints),
			StartTime:   *start,
			Time:        timeVal,
			Temporality: toWasiTemporality(v.Temporality),
		})
	case metricdata.Histogram[int64]:
		start, timeVal := extractTimestamps(v.DataPoints)
		return wasi_otel_metrics.MakeMetricDataS64Histogram(wasi_otel_metrics.Histogram{
			DataPoints:  toWasiHistogramDataPoint(v.DataPoints),
			StartTime:   *start,
			Time:        timeVal,
			Temporality: toWasiTemporality(v.Temporality),
		})
	case metricdata.Histogram[float64]:
		start, timeVal := extractTimestamps(v.DataPoints)
		return wasi_otel_metrics.MakeMetricDataF64Histogram(wasi_otel_metrics.Histogram{
			DataPoints:  toWasiHistogramDataPoint(v.DataPoints),
			StartTime:   *start,
			Time:        timeVal,
			Temporality: toWasiTemporality(v.Temporality),
		})
	case metricdata.Summary:
		panic("The metricdata.Summary metric type is not implemented")
	default:
		panic("unimplemented type")
	}
}

func toWasiGaugeDataPoint[T float64 | int64](dataPoints []metricdata.DataPoint[T]) []wasi_otel_metrics.GaugeDataPoint {
	result := make([]wasi_otel_metrics.GaugeDataPoint, len(dataPoints))
	for i, dp := range dataPoints {
		result[i] = wasi_otel_metrics.GaugeDataPoint{
			Attributes: types.ToWasiAttributes(dp.Attributes.ToSlice()),
			Value:      toWasiMetricNumber(dp.Value),
			Exemplars:  toWasiExemplar(dp.Exemplars),
		}
	}

	return result
}

func toWasiSumDataPoint[T float64 | int64](dataPoints []metricdata.DataPoint[T]) []wasi_otel_metrics.SumDataPoint {
	result := make([]wasi_otel_metrics.SumDataPoint, len(dataPoints))
	for i, dp := range dataPoints {
		result[i] = wasi_otel_metrics.SumDataPoint{
			Attributes: types.ToWasiAttributes(dp.Attributes.ToSlice()),
			Value:      toWasiMetricNumber(dp.Value),
			Exemplars:  toWasiExemplar(dp.Exemplars),
		}
	}

	return result
}

func toWasiHistogramDataPoint[T float64 | int64](dataPoints []metricdata.HistogramDataPoint[T]) []wasi_otel_metrics.HistogramDataPoint {
	result := make([]wasi_otel_metrics.HistogramDataPoint, len(dataPoints))
	for i, dp := range dataPoints {
		result[i] = wasi_otel_metrics.HistogramDataPoint{
			Attributes:   types.ToWasiAttributes(dp.Attributes.ToSlice()),
			Count:        dp.Count,
			Bounds:       dp.Bounds,
			BucketCounts: dp.BucketCounts,
			Min:          toWasiOptMetricNumber(dp.Min),
			Max:          toWasiOptMetricNumber(dp.Max),
			Sum:          toWasiMetricNumber(dp.Sum),
			Exemplars:    toWasiExemplar(dp.Exemplars),
		}
	}

	return result
}

func toWasiExponentialHistogram[T float64 | int64](dataPoints []metricdata.ExponentialHistogramDataPoint[T]) []wasi_otel_metrics.ExponentialHistogramDataPoint {
	result := make([]wasi_otel_metrics.ExponentialHistogramDataPoint, len(dataPoints))
	for i, dp := range dataPoints {
		result[i] = wasi_otel_metrics.ExponentialHistogramDataPoint{
			Attributes: types.ToWasiAttributes(dp.Attributes.ToSlice()),
			Count:      dp.Count,
			Min:        toWasiOptMetricNumber(dp.Min),
			Max:        toWasiOptMetricNumber(dp.Max),
			Sum:        toWasiMetricNumber(dp.Sum),
			Scale:      int8(dp.Scale),
			ZeroCount:  dp.ZeroCount,
			PositiveBucket: wasi_otel_metrics.ExponentialBucket{
				Offset: dp.PositiveBucket.Offset,
				Counts: dp.PositiveBucket.Counts,
			},
			NegativeBucket: wasi_otel_metrics.ExponentialBucket{
				Offset: dp.NegativeBucket.Offset,
				Counts: dp.NegativeBucket.Counts,
			},
			ZeroThreshold: dp.ZeroThreshold,
			Exemplars:     toWasiExemplar(dp.Exemplars),
		}
	}

	return result
}

func toWasiOptMetricNumber[T float64 | int64](n metricdata.Extrema[T]) wit_types.Option[wasi_otel_metrics.MetricNumber] {
	num, exists := n.Value()
	if exists {
		return wit_types.Some(toWasiMetricNumber(num))
	} else {
		return wit_types.None[wasi_otel_metrics.MetricNumber]()
	}
}

func toWasiTemporality(t metricdata.Temporality) wasi_otel_metrics.Temporality {
	switch t {
	case metricdata.CumulativeTemporality:
		return wasi_otel_metrics.TemporalityCumulative
	case metricdata.DeltaTemporality:
		return wasi_otel_metrics.TemporalityDelta
	default:
		return wasi_otel_metrics.TemporalityCumulative
	}
}

func toWasiMetricNumber[T float64 | int64](n T) wasi_otel_metrics.MetricNumber {
	switch v := any(n).(type) {
	case int64:
		return wasi_otel_metrics.MakeMetricNumberS64(v)
	case float64:
		return wasi_otel_metrics.MakeMetricNumberF64(v)
	default:
		panic("unsupported type")
	}
}

func toWasiExemplar[T float64 | int64](exemplars []metricdata.Exemplar[T]) []wasi_otel_metrics.Exemplar {
	result := make([]wasi_otel_metrics.Exemplar, len(exemplars))
	for i, e := range exemplars {
		result[i] = wasi_otel_metrics.Exemplar{
			FilteredAttributes: types.ToWasiAttributes(e.FilteredAttributes),
			Time:               types.ToWasiTime(e.Time),
			Value:              toWasiMetricNumber(e.Value),
			SpanId:             string(e.SpanID),
			TraceId:            string(e.TraceID),
		}
	}

	return result
}

// timestampProvider is a constraint for types that have StartTime and Time fields
type timestampProvider interface {
	metricdata.DataPoint[int64] | metricdata.DataPoint[float64] |
		metricdata.HistogramDataPoint[int64] | metricdata.HistogramDataPoint[float64] |
		metricdata.ExponentialHistogramDataPoint[int64] | metricdata.ExponentialHistogramDataPoint[float64]
}

// extractTimestamps extracts StartTime and Time from the first data point in a slice
func extractTimestamps[T timestampProvider](dataPoints []T) (startTime *wasi_clocks_wall_clock.Datetime, timeRecorded wasi_clocks_wall_clock.Datetime) {
	if len(dataPoints) == 0 {
		return nil, types.ToWasiTime(time.Now())
	}

	// Use type assertion to access the fields
	var start, timeVal time.Time
	switch dp := any(dataPoints[0]).(type) {
	case metricdata.DataPoint[int64]:
		start, timeVal = dp.StartTime, dp.Time
	case metricdata.DataPoint[float64]:
		start, timeVal = dp.StartTime, dp.Time
	case metricdata.HistogramDataPoint[int64]:
		start, timeVal = dp.StartTime, dp.Time
	case metricdata.HistogramDataPoint[float64]:
		start, timeVal = dp.StartTime, dp.Time
	case metricdata.ExponentialHistogramDataPoint[int64]:
		start, timeVal = dp.StartTime, dp.Time
	case metricdata.ExponentialHistogramDataPoint[float64]:
		start, timeVal = dp.StartTime, dp.Time
	default:
		panic("unsupported types")
	}

	resultStart := types.ToWasiTime(start)

	return &resultStart, types.ToWasiTime(timeVal)
}

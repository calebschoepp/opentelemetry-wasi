import {
  AggregationTemporality,
  AggregationType,
  ExponentialHistogramMetricData,
  GaugeMetricData,
  HistogramMetricData,
  MetricData,
  MetricReader,
  ResourceMetrics,
  ScopeMetrics,
  SumMetricData,
  DataPointType,
} from '@opentelemetry/sdk-metrics';
import {
  export as exportToWasiHost,
  Metric as WasiMetric,
  ResourceMetrics as WasiResourceMetrics,
  ScopeMetrics as WasiScopeMetrics,
  SumDataPoint as WasiSumDataPoint,
  GaugeDataPoint as WasiGaugeDataPoint,
  HistogramDataPoint as WasiHistogramDataPoint,
  ExponentialHistogramDataPoint as WasiExponentialHistogramDataPoint,
  Temporality as WasiTemporality,
  MetricNumber as WasiMetricNumber,
} from 'wasi:otel/metrics@0.2.0-rc.2';
import { diag, ValueType, HrTime } from '@opentelemetry/api';
import { dateTimeToWasi, attributesToWasi } from './types';

export class WasiMetricExporter extends MetricReader {
  constructor() {
    super({
      aggregationSelector: (_instrumentType) => {
        return {
          type: AggregationType.DEFAULT,
        };
      },
      aggregationTemporalitySelector: (_instrumentType) =>
        AggregationTemporality.CUMULATIVE,
    });
  }

  protected override async onForceFlush(): Promise<void> {
    // no-op
  }

  protected override async onShutdown(): Promise<void> {
    // no-op
  }

  /**
   * Exports metric data to a compatible host or component.
   */
  public async export(): Promise<void> {
    const { resourceMetrics, errors } = await this.collect();
    if (errors.length) {
      diag.error('WasiMetricExporter: metrics collection errors', ...errors);
    }
    exportToWasiHost(resourceMetricsToWasi(resourceMetrics));
  }
}

/**
 * Converts OpenTelemetry ResourceMetrics to WASI ResourceMetrics
 */
function resourceMetricsToWasi(rm: ResourceMetrics): WasiResourceMetrics {
  return {
    resource: {
      attributes: attributesToWasi(rm.resource.attributes),
      schemaUrl: rm.resource.schemaUrl,
    },
    scopeMetrics: scopeMetricsToWasi(rm.scopeMetrics),
  };
}

/**
 * Converts OpenTelemetry ScopeMetrics to WASI ScopeMetrics
 */
function scopeMetricsToWasi(scopeMetrics: ScopeMetrics[]): WasiScopeMetrics[] {
  return scopeMetrics.map((sm) => ({
    scope: {
      name: sm.scope.name,
      attributes: [],
    },
    metrics: metricsToWasi(sm.metrics),
  }));
}

/**
 * Converts an array of OpenTelemetry MetricData to an array of WASI Metric
 */
function metricsToWasi(metrics: MetricData[]): WasiMetric[] {
  return metrics.map((m) => {
    switch (m.dataPointType) {
      case DataPointType.SUM:
        return sumToWasi(m);
      case DataPointType.GAUGE:
        return gaugeToWasi(m);
      case DataPointType.HISTOGRAM:
        return histogramToWasi(m);
      case DataPointType.EXPONENTIAL_HISTOGRAM:
        return exponentialHistogramToWasi(m);
      default:
        throw new Error(`Unknown data point type`);
    }
  });
}

/**
 * Converts OpenTelemetry SumMetricData to WASI Metric
 */
function sumToWasi(data: SumMetricData): WasiMetric {
  const isF64 = data.descriptor.valueType === ValueType.DOUBLE;
  const dataPoints: WasiSumDataPoint[] = data.dataPoints.map((dp) => ({
    attributes: attributesToWasi(dp.attributes),
    value: numberToWasi(isF64, dp.value),
    // Exemplars are defined in the package, but are not exported.
    // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/sdk-metrics/src/index.ts
    exemplars: [],
  }));

  const sumData = {
    dataPoints,
    ...getTimeRange(data.dataPoints),
    isMonotonic: data.isMonotonic,
    temporality: temporalityToWasi(data.aggregationTemporality),
  };

  return {
    description: data.descriptor.description,
    name: data.descriptor.name,
    unit: data.descriptor.unit,
    data: isF64
      ? { tag: 'f64-sum', val: sumData }
      : { tag: 's64-sum', val: sumData },
  };
}

/**
 * Converts OpenTelemetry GaugeMetricData to WASI Metric
 */
function gaugeToWasi(data: GaugeMetricData): WasiMetric {
  const isF64 = data.descriptor.valueType === ValueType.DOUBLE;

  const dataPoints: WasiGaugeDataPoint[] = data.dataPoints.map((dp) => ({
    attributes: attributesToWasi(dp.attributes),
    value: numberToWasi(isF64, dp.value),
    // Exemplars are defined in the package, but are not exported.
    // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/sdk-metrics/src/index.ts
    exemplars: [],
  }));

  const gaugeData = {
    dataPoints,
    ...getTimeRange(data.dataPoints),
  };

  return {
    description: data.descriptor.description,
    name: data.descriptor.name,
    unit: data.descriptor.unit,
    data: isF64
      ? { tag: 'f64-gauge', val: gaugeData }
      : { tag: 's64-gauge', val: gaugeData },
  };
}

/**
 * Converts OpenTelemetry HistogramMetricData to WASI Metric
 */
function histogramToWasi(data: HistogramMetricData): WasiMetric {
  const isF64 = data.descriptor.valueType === ValueType.DOUBLE;
  const dataPoints: WasiHistogramDataPoint[] = data.dataPoints.map((dp) => ({
    attributes: attributesToWasi(dp.attributes),
    count: BigInt(dp.value.count),
    max: numberToWasi(isF64, dp.value.max),
    min: numberToWasi(isF64, dp.value.min),
    sum: numberToWasi(isF64, dp.value.sum),
    bounds: new Float64Array(dp.value.buckets.boundaries),
    bucketCounts: new BigUint64Array(dp.value.buckets.counts.map(BigInt)),
    // Exemplars are defined in the package, but are not exported.
    // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/sdk-metrics/src/index.ts
    exemplars: [],
  }));

  const histData = {
    dataPoints,
    ...getTimeRange(data.dataPoints),
    temporality: temporalityToWasi(data.aggregationTemporality),
  };

  return {
    description: data.descriptor.description,
    name: data.descriptor.name,
    unit: data.descriptor.unit,
    data: isF64
      ? { tag: 'f64-histogram', val: histData }
      : { tag: 's64-histogram', val: histData },
  };
}

/**
 * Converts OpenTelemetry ExponentialHistogramMetricData to WASI Metric
 */
function exponentialHistogramToWasi(
  data: ExponentialHistogramMetricData
): WasiMetric {
  const isF64 = data.descriptor.valueType === ValueType.DOUBLE;
  const dataPoints: WasiExponentialHistogramDataPoint[] = data.dataPoints.map(
    (dp) => ({
      attributes: attributesToWasi(dp.attributes),
      count: BigInt(dp.value.count),
      max: numberToWasi(isF64, dp.value.max),
      min: numberToWasi(isF64, dp.value.min),
      sum: numberToWasi(isF64, dp.value.sum),
      scale: dp.value.scale,
      zeroCount: BigInt(dp.value.zeroCount),
      positiveBucket: {
        offset: dp.value.positive.offset,
        counts: new BigUint64Array(dp.value.positive.bucketCounts.map(BigInt)),
      },
      negativeBucket: {
        offset: dp.value.negative.offset,
        counts: new BigUint64Array(dp.value.negative.bucketCounts.map(BigInt)),
      },
      // According to the spec, the zero threshold defaults to 0.
      // See https://opentelemetry.io/docs/specs/otel/metrics/data-model/#zero-count-and-zero-threshold
      zeroThreshold: 0,
      // Exemplars are defined in the package, but are not exported.
      // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/sdk-metrics/src/index.ts
      exemplars: [],
    })
  );

  const histData = {
    dataPoints,
    ...getTimeRange(data.dataPoints),
    temporality: temporalityToWasi(data.aggregationTemporality),
  };

  return {
    description: data.descriptor.description,
    name: data.descriptor.name,
    unit: data.descriptor.unit,
    data: isF64
      ? { tag: 'f64-exponential-histogram', val: histData }
      : { tag: 's64-exponential-histogram', val: histData },
  };
}

function temporalityToWasi(t: AggregationTemporality): WasiTemporality {
  return (
    t === AggregationTemporality.DELTA ? 'delta' : 'cumulative'
  ) as WasiTemporality;
}

/**
 * Converts a number to WASI MetricNumber
 */
function numberToWasi(
  isF64: boolean,
  value: number | undefined
): WasiMetricNumber {
  const v = value ?? 0;
  return isF64 ? { tag: 'f64', val: v } : { tag: 's64', val: BigInt(v) };
}

/**
 * Parses the start and end times from a list of datapoints.
 * @returns An object containing the start and end times.
 */
function getTimeRange(dataPoints: { startTime: HrTime; endTime: HrTime }[]) {
  const defaultTime = dateTimeToWasi([0, 0]); // TODO: I wonder if there's a better way to handle this
  return {
    startTime: dataPoints[0]?.startTime
      ? dateTimeToWasi(dataPoints[0].startTime)
      : defaultTime,
    time: dataPoints[0]?.endTime
      ? dateTimeToWasi(dataPoints[0].endTime)
      : defaultTime,
  };
}

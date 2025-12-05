import {
  TraceFlags as WasiTraceFlags,
  TraceState as WasiTraceState,
  SpanContext as WasiSpanContext,
  SpanKind as WasiSpanKind,
  Datetime,
  SpanData as WasiSpanData,
} from "wasi:otel/tracing@0.2.0-draft";
import {
  ResourceMetrics as WasiResourceMetrics,
  ScopeMetrics as WasiScopeMetrics,
  Metric as WasiMetric,
  Temporality as WasiTemporality,
  SumDataPoint as WasiSumDataPoint,
  GaugeDataPoint as WasiGaugeDataPoint,
  HistogramDataPoint as WasiHistogramDataPoint,
  ExponentialHistogramDataPoint as WasiExponentialHistogramDataPoint,
  MetricNumber as WasiMetricNumber,
} from "wasi:otel/metrics@0.2.0-draft";
import {
  Value,
  KeyValue,
} from "wasi:otel/types@0.2.0-draft";
import {
    Attributes,
    AttributeValue,
    createTraceState,
    HrTime,
    SpanContext,
    TraceFlags,
    TraceState,
    ValueType,
} from "@opentelemetry/api";
import {
  DataPointType,
  ExponentialHistogramMetricData,
  GaugeMetricData,
  HistogramMetricData,
  MetricData,
  ResourceMetrics,
  ScopeMetrics,
  SumMetricData
} from "@opentelemetry/sdk-metrics";
import { ReadableSpan } from "@opentelemetry/sdk-trace-base";

/**
 * Converts OpenTelemetry SpanContext to WASI SpanContext
 */
export function spanContextToWasi(ctx: SpanContext): WasiSpanContext {
  return {
    traceId: ctx.traceId,
    spanId: ctx.spanId,
    traceFlags: traceFlagsToWasi(ctx.traceFlags),
    isRemote: ctx.isRemote || false,
    traceState: traceStateToWasi(ctx.traceState),
  }
}

/**
 * Converts WASI TraceState to OpenTelemetry TraceState
 */
function wasiToTraceState (wts: WasiTraceState): TraceState {
  let traceState = createTraceState();
  for (const [key, value] of wts) {
      traceState = traceState.set(key, value);
  }

  return traceState;
}

/**
 * Converts WASI SpanContext to OpenTelemetry SpanContext
 */
export function wasiToSpanContext(ctx: WasiSpanContext): SpanContext {
    return {
        traceId: ctx.traceId,
        spanId: ctx.spanId,
        traceFlags: ctx.traceFlags.sampled ? TraceFlags.SAMPLED : TraceFlags.NONE,
        isRemote: ctx.isRemote || false,
        traceState: wasiToTraceState(ctx.traceState),
    }
}

/**
 * Converts OpenTelemetry Attributes to WASI Attributes
 */
export function attributesToWasi(attrs: Attributes | undefined): KeyValue[] {
  if (attrs === undefined) {
    return [];
  }
  return Object.entries(attrs).map(([key, value]) => ({
    key,
    value: attributeValueToWasi(value),
  }));
}

/**
 * Converts OpenTelemetry AttributeValue to WASI AttributeValue
 */
function attributeValueToWasi(value: AttributeValue | undefined): Value {
  if (typeof value === 'string') {
    return { tag: 'string', val: value };
  } else if (typeof value === 'boolean') {
    return { tag: 'bool', val: value };
  } else if (typeof value === 'number') {
    return Number.isInteger(value)
      ? { tag: 's64', val: BigInt(value) }
      : { tag: 'f64', val: value };
  } else if (Array.isArray(value)) {
    const filtered = value.filter(v => v != null);
    if (filtered.length === 0) {
      return { tag: 'string-array', val: [] };
    }
    const firstType = typeof filtered[0];
    if (firstType === 'string') {
      return { tag: 'string-array', val: filtered as string[] };
    } else if (firstType === 'boolean') {
      return { tag: 'bool-array', val: filtered as boolean[] };
    } else if (firstType === 'number') {
      const numbers = filtered as number[];
      if (numbers.every(Number.isInteger)) {
        const bigIntArray = new BigInt64Array(numbers.length);
        numbers.forEach((n, i) => bigIntArray[i] = BigInt(n));
        return { tag: 's64-array', val: bigIntArray };
      } else {
        return { tag: 'f64-array', val: new Float64Array(numbers) };
      }
    }
  }
  // Fallback
  return { tag: 'string', val: String(value) };
}

/**
 * Converts OpenTelemetry TraceFlags to WASI TraceFlags
 */
function traceFlagsToWasi(flags: TraceFlags): WasiTraceFlags {
  const SAMPLED = 0x01;
  return (flags & SAMPLED) === SAMPLED ? {sampled: true} : {sampled: false};
}

/**
 * Converts OpenTelemetry TraceState to WASI TraceState
 */
function traceStateToWasi(value: TraceState | undefined): WasiTraceState {
  if (value == undefined) {
    return [];
  }
  return value.serialize().split(',').map(entry => {
     // This ensures that a pattern like "foo=bar=baz" is split into
     // Key("foo"), Value("bar=baz")
      const [key, ...rest] = entry.split('=');
      const value = rest.join('=');
      return [key, value];
  });
}

/**
 * Converts OpenTelemetry HrTime to WASI Datetime
 */
export function dateTimeToWasi(time: HrTime): Datetime {
  return {
    seconds: BigInt(time[0]),
    nanoseconds: time[1],
  };
}

/**
 * Converts OpenTelemetry ReadableSpan to WASI SpanData
 */
export function spanDataToWasi(span: ReadableSpan): WasiSpanData {
  return {
      name: span.name,
      startTime: dateTimeToWasi(span.startTime),
      spanContext: spanContextToWasi(span.spanContext()),
      parentSpanId: span.parentSpanContext?.spanId || "",
      spanKind: ["internal", "server", "client", "producer", "consumer"][span.kind] as WasiSpanKind,
      endTime: dateTimeToWasi(span.endTime),
      attributes: attributesToWasi(span.attributes),
      events: span.events.map(e => ({
        name: e.name,
        time: dateTimeToWasi(e.time),
        attributes: attributesToWasi(e.attributes),
      })),
      links: span.links.map(link => ({
        spanContext: spanContextToWasi(link.context),
        attributes: attributesToWasi(link.attributes),
      })),
      status: span.status.code === 0 ? { tag: 'unset' } :
        span.status.code === 1 ? { tag: 'ok' } :
        { tag: 'error', val: span.status.message || "" },
      instrumentationScope: {
        name: span.instrumentationScope.name,
        version: span.instrumentationScope.version,
        schemaUrl: span.instrumentationScope.schemaUrl,
        // Although other SDKs use the InstrumentationScope.attributes field, the `opentelemetry-js` SDK does not.
        // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/opentelemetry-core/src/common/types.ts#L47
        attributes: [],
      },
      droppedAttributes: span.droppedAttributesCount,
      droppedEvents: span.droppedEventsCount,
      droppedLinks: span.droppedLinksCount,
    };
}

/**
 * Converts OpenTelemetry ResourceMetrics to WASI ResourceMetrics
 */
export function resourceMetricsToWasi(rm: ResourceMetrics): WasiResourceMetrics {
  return {
    resource: {
      attributes: attributesToWasi(rm.resource.attributes),
      schemaUrl: rm.resource.schemaUrl,

    },
    scopeMetrics: scopeMetricsToWasi(rm.scopeMetrics),
  }
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
        throw new Error(`Unknown data point type: ${typeof m}`);
    }
  });
}

/**
 * Converts OpenTelemetry SumMetricData to WASI Metric
 */
function sumToWasi(data: SumMetricData): WasiMetric {
  const isF64 = data.descriptor.valueType === ValueType.DOUBLE;
  const defaultTime = dateTimeToWasi([0, 0]); // TODO: I wonder if there's a better way to handle this
  const dataPoints: WasiSumDataPoint[] = data.dataPoints.map((dp) => ({
    attributes: attributesToWasi(dp.attributes),
    value: numberToWasi(isF64, dp.value),
    // Exemplars are defined in the package, but are not exported.
    // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/sdk-metrics/src/index.ts
    exemplars: [],
  }));

  const sumData = {
    dataPoints,
    startTime: data.dataPoints[0]?.startTime
      ? dateTimeToWasi(data.dataPoints[0].startTime)
      : defaultTime,
    time: data.dataPoints[0]?.endTime
      ? dateTimeToWasi(data.dataPoints[0].endTime)
      : defaultTime,
    isMonotonic: data.isMonotonic,
    temporality: (data.aggregationTemporality === 0 ? "delta" : "cumulative") as WasiTemporality,
  };

  return {
    description: data.descriptor.description,
    name: data.descriptor.name,
    unit: data.descriptor.unit,
    data: isF64
      ? { tag: "f64-sum", val: sumData }
      : { tag: "s64-sum", val: sumData },
  };
}

/**
 * Converts OpenTelemetry GaugeMetricData to WASI Metric
 */
function gaugeToWasi(data: GaugeMetricData): WasiMetric {
  const isF64 = data.descriptor.valueType === ValueType.DOUBLE;
  const defaultTime = dateTimeToWasi([0, 0]); // TODO: I wonder if there's a better way to handle this

  const dataPoints: WasiGaugeDataPoint[] = data.dataPoints.map((dp) => ({
    attributes: attributesToWasi(dp.attributes),
    value: numberToWasi(isF64, dp.value),
    // Exemplars are defined in the package, but are not exported.
    // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/sdk-metrics/src/index.ts
    exemplars: [],
  }));

  const gaugeData = {
    dataPoints: dataPoints,
    startTime: data.dataPoints[0]?.startTime
    ? dateTimeToWasi(data.dataPoints[0].startTime)
    : defaultTime,
  time: data.dataPoints[0]?.endTime
    ? dateTimeToWasi(data.dataPoints[0].endTime)
    : defaultTime,
  }

  return {
    description: data.descriptor.description,
    name: data.descriptor.name,
    unit: data.descriptor.unit,
    data: isF64
      ? { tag: "f64-gauge", val: gaugeData }
      : { tag: "s64-gauge", val: gaugeData },
  };
}

/**
 * Converts OpenTelemetry HistogramMetricData to WASI Metric
 */
function histogramToWasi(data: HistogramMetricData): WasiMetric {
  const defaultTime = dateTimeToWasi([0, 0]); // TODO: I wonder if there's a better way to handle this
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
    startTime: data.dataPoints[0]?.startTime
      ? dateTimeToWasi(data.dataPoints[0].startTime)
      : defaultTime,
    time: data.dataPoints[0]?.endTime
      ? dateTimeToWasi(data.dataPoints[0].endTime)
      : defaultTime,
    temporality: (data.aggregationTemporality === 0 ? "delta" : "cumulative") as WasiTemporality,
  };

  return {
    description: data.descriptor.description,
    name: data.descriptor.name,
    unit: data.descriptor.unit,
    data: isF64
      ? { tag: "f64-histogram", val: histData }
      : { tag: "s64-histogram", val: histData },
  };
}

/**
 * Converts OpenTelemetry ExponentailHistogramMetricData to WASI Metric
 */
function exponentialHistogramToWasi(data: ExponentialHistogramMetricData): WasiMetric {
  const defaultTime = dateTimeToWasi([0, 0]); // TODO: I wonder if there's a better way to handle this
  const isF64 = data.descriptor.valueType === ValueType.DOUBLE;
  const dataPoints: WasiExponentialHistogramDataPoint[] = data.dataPoints.map((dp) => ({
    attributes: attributesToWasi(dp.attributes),
    count: BigInt(dp.value.count),
    max: numberToWasi(isF64, dp.value.max),
    min: numberToWasi(isF64, dp.value.min),
    sum: numberToWasi(isF64, dp.value.sum),
    scale: dp.value.scale,
    zeroCount: BigInt(dp.value.zeroCount),
    positiveBucket: {offset: dp.value.positive.offset, counts: new BigUint64Array(dp.value.positive.bucketCounts.map(BigInt))},
    negativeBucket: {offset: dp.value.negative.offset, counts: new BigUint64Array(dp.value.negative.bucketCounts.map(BigInt))},
    // According to the spec, the zero threshold defaults to 0.
    // See https://opentelemetry.io/docs/specs/otel/metrics/data-model/#zero-count-and-zero-threshold
    zeroThreshold: 0,
    // Exemplars are defined in the package, but are not exported.
    // See https://github.com/open-telemetry/opentelemetry-js/blob/main/packages/sdk-metrics/src/index.ts
    exemplars: [],
  }));

  const histData = {
    dataPoints,
    startTime: data.dataPoints[0]?.startTime
      ? dateTimeToWasi(data.dataPoints[0].startTime)
      : defaultTime,
    time: data.dataPoints[0]?.endTime
      ? dateTimeToWasi(data.dataPoints[0].endTime)
      : defaultTime,
    temporality: (data.aggregationTemporality === 0 ? "delta" : "cumulative") as WasiTemporality,
  };

  return {
    description: data.descriptor.description,
    name: data.descriptor.name,
    unit: data.descriptor.unit,
    data: isF64
      ? { tag: "f64-exponential-histogram", val: histData }
      : { tag: "s64-exponential-histogram", val: histData },
  };
}

/**
 * Converts a number to WASI MetricNumber
 */
function numberToWasi(isF64: boolean, value: number | undefined): WasiMetricNumber {
  return isF64
    ? {tag: "f64", val: value ? value : 0}
      : {tag: "s64", val: BigInt(value ? value : 0)}
}

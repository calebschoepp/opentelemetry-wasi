/// <reference path="./wasi-clocks-monotonic-clock.d.ts" />
/// <reference path="./wasi-clocks-wall-clock.d.ts" />
/// <reference path="./wasi-otel-tracing.d.ts" />
/// <reference path="./wasi-otel-types.d.ts" />
declare module 'wasi:otel/metrics@0.2.0-draft' {
  /**
   * `collect` gathers all metric data related to a Reader from the SDK
   */
  export function collect(metrics: ResourceMetrics): void;
  export type Datetime = import('wasi:clocks/wall-clock@0.2.0').Datetime;
  export type Duration = import('wasi:clocks/monotonic-clock@0.2.0').Duration;
  export type KeyValue = import('wasi:otel/types@0.2.0-draft').KeyValue;
  export type InstrumentationScope = import('wasi:otel/types@0.2.0-draft').InstrumentationScope;
  export type SpanId = import('wasi:otel/tracing@0.2.0-draft').SpanId;
  export type TraceId = import('wasi:otel/tracing@0.2.0-draft').TraceId;
  /**
   * Inner structure of `resource` holding the actual data.
   */
  export interface ResourceInner {
    attributes: Array<KeyValue>,
    schemaUrl?: string,
  }
  /**
   * An immutable representation of the entity producing telemetry as attributes.
   */
  export interface Resource {
    inner: ResourceInner,
  }
  /**
   * A set of bucket counts, encoded in a contiguous array of counts.
   */
  export interface ExponentialBucket {
    /**
     * The bucket index of the first entry in the `counts` list.
     */
    offset: number,
    /**
     * A list where `counts[i]` carries the count of the bucket at index `offset + i`.
     * 
     * `counts[i]` is the count of values greater than base^(offset+i) and less than
     * or equal to base^(offset+i+1).
     */
    counts: BigUint64Array,
  }
  /**
   * Defines the window that an aggregation was calculated over.
   * # Variants
   * 
   * ## `"cumulative"`
   * 
   * A measurement interval that continues to expand forward in time from a
   * starting point.
   * 
   * New measurements are added to all previous measurements since a start time.
   * 
   * This is the default temporality
   * ## `"delta"`
   * 
   * A measurement interval that resets each cycle.
   * 
   * Measurements from one cycle are recorded independently, measurements from
   * other cycles do not affect them.
   * ## `"low-memory"`
   * 
   * Configures Synchronous Counter and Histogram instruments to use
   * Delta aggregation temporality, which allows them to shed memory
   * following a cardinality explosion, thus use less memory.
   */
  export type Temporality = 'cumulative' | 'delta' | 'low-memory';
  /**
   * The number types available for use with the OpenTelemetry SDKs.
   * 
   * This makes it easier to use generics when converting to and from WASI data types.
   */
  export type MetricNumber = MetricNumberF64 | MetricNumberS64 | MetricNumberU64;
  export interface MetricNumberF64 {
    tag: 'f64',
    val: number,
  }
  export interface MetricNumberS64 {
    tag: 's64',
    val: bigint,
  }
  export interface MetricNumberU64 {
    tag: 'u64',
    val: bigint,
  }
  /**
   * A measurement sampled from a time series providing a typical example.
   */
  export interface Exemplar {
    /**
     * The attributes recorded with the measurement but filtered out of the
     * time series' aggregated data.
     */
    filteredAttributes: Array<KeyValue>,
    /**
     * The time when the measurement was recorded.
     */
    time: Datetime,
    /**
     * The measured value.
     */
    value: MetricNumber,
    /**
     * The ID of the span that was active during the measurement.
     * 
     * If no span was active or the span was not sampled this will be empty.
     */
    spanId: SpanId,
    /**
     * The ID of the trace the active span belonged to during the measurement.
     * 
     * If no span was active or the span was not sampled this will be empty.
     */
    traceId: TraceId,
  }
  /**
   * `gauge-data-point` is a single data point in a time series.
   */
  export interface GaugeDataPoint {
    /**
     * `attributes` is the set of key value pairs that uniquely identify the
     * time series.
     */
    attributes: Array<KeyValue>,
    /**
     * The value of this data point.
     */
    value: MetricNumber,
    /**
     * The sampled `exemplar`s collected during the time series.
     */
    exemplars: Array<Exemplar>,
  }
  /**
   * A measurement of the current value of an instrument.
   */
  export interface Gauge {
    /**
     * Represents individual aggregated measurements with unique attributes.
     */
    dataPoints: Array<GaugeDataPoint>,
    /**
     * The time when the time series was started.
     */
    startTime?: Datetime,
    /**
     * The time when the time series was recorded.
     */
    time: Datetime,
  }
  /**
   * `sum-data-point` is a single data point in a time series.
   */
  export interface SumDataPoint {
    /**
     * `attributes` is the set of key value pairs that uniquely identify the
     * time series.
     */
    attributes: Array<KeyValue>,
    /**
     * The value of this data point.
     */
    value: MetricNumber,
    /**
     * The sampled `exemplar`s collected during the time series.
     */
    exemplars: Array<Exemplar>,
  }
  /**
   * Represents the sum of all measurements of values from an instrument.
   */
  export interface Sum {
    /**
     * Represents individual aggregated measurements with unique attributes.
     */
    dataPoints: Array<SumDataPoint>,
    /**
     * The time when the time series was started.
     */
    startTime: Datetime,
    /**
     * The time when the time series was recorded.
     */
    time: Datetime,
    /**
     * Describes if the aggregation is reported as the change from the last report
     * time, or the cumulative changes since a fixed start time.
     */
    temporality: Temporality,
    /**
     * Whether this aggregation only increases or decreases.
     */
    isMonotonic: boolean,
  }
  /**
   * A single histogram data point in a time series.
   */
  export interface HistogramDataPoint {
    /**
     * The set of key value pairs that uniquely identify the time series.
     */
    attributes: Array<KeyValue>,
    /**
     * The number of updates this histogram has been calculated with.
     */
    count: bigint,
    /**
     * The upper bounds of the buckets of the histogram.
     */
    bounds: Float64Array,
    /**
     * The count of each of the buckets.
     */
    bucketCounts: BigUint64Array,
    /**
     * The minimum value recorded.
     */
    min?: MetricNumber,
    /**
     * The maximum value recorded.
     */
    max?: MetricNumber,
    /**
     * The sum of the values recorded
     */
    sum: MetricNumber,
    /**
     * The sampled `exemplar`s collected during the time series.
     */
    exemplars: Array<Exemplar>,
  }
  /**
   * Represents the histogram of all measurements of values from an instrument.
   */
  export interface Histogram {
    /**
     * Individual aggregated measurements with unique attributes.
     */
    dataPoints: Array<HistogramDataPoint>,
    /**
     * The time when the time series was started.
     */
    startTime: Datetime,
    /**
     * The time when the time series was recorded.
     */
    time: Datetime,
    /**
     * Describes if the aggregation is reported as the change from the last report
     * time, or the cumulative changes since a fixed start time.
     */
    temporality: Temporality,
  }
  /**
   * A single exponential histogram data point in a time series.
   */
  export interface ExponentialHistogramDataPoint {
    /**
     * The set of key value pairs that uniquely identify the time series.
     */
    attributes: Array<KeyValue>,
    /**
     * The number of updates this histogram has been calculated with.
     */
    count: bigint,
    /**
     * TODO: check that u64 is an acceptable replacement for usize
     * The minimum value recorded.
     */
    min?: MetricNumber,
    /**
     * The maximum value recorded.
     */
    max?: MetricNumber,
    /**
     * The maximum value recorded.
     */
    sum: MetricNumber,
    /**
     * Describes the resolution of the histogram.
     * 
     * Boundaries are located at powers of the base, where:
     * 
     *   base = 2 ^ (2 ^ -scale)
     */
    scale: number,
    /**
     * The number of values whose absolute value is less than or equal to
     * `zero_threshold`.
     * 
     * When `zero_threshold` is `0`, this is the number of values that cannot be
     * expressed using the standard exponential formula as well as values that have
     * been rounded to zero.
     */
    zeroCount: bigint,
    /**
     * The range of positive value bucket counts.
     */
    positiveBucket: ExponentialBucket,
    /**
     * The range of negative value bucket counts.
     */
    negativeBucket: ExponentialBucket,
    /**
     * The width of the zero region.
     * 
     * Where the zero region is defined as the closed interval
     * [-zero_threshold, zero_threshold].
     */
    zeroThreshold: number,
    /**
     * The sampled exemplars collected during the time series.
     */
    exemplars: Array<Exemplar>,
  }
  /**
   * The histogram of all measurements of values from an instrument.
   */
  export interface ExponentialHistogram {
    /**
     * The individual aggregated measurements with unique attributes.
     */
    dataPoints: Array<ExponentialHistogramDataPoint>,
    /**
     * When the time series was started.
     */
    startTime: Datetime,
    /**
     * The time when the time series was recorded.
     */
    time: Datetime,
    /**
     * Describes if the aggregation is reported as the change from the last report
     * time, or the cumulative changes since a fixed start time.
     */
    temporality: Temporality,
  }
  /**
   * Metric data for all types.
   */
  export type MetricData = MetricDataGauge | MetricDataSum | MetricDataHistogram | MetricDataExponentialHistogram;
  /**
   * Metric data for `gauge`.
   */
  export interface MetricDataGauge {
    tag: 'gauge',
    val: Gauge,
  }
  /**
   * Metric data for `sum`.
   */
  export interface MetricDataSum {
    tag: 'sum',
    val: Sum,
  }
  /**
   * Metric data for `histogram`.
   */
  export interface MetricDataHistogram {
    tag: 'histogram',
    val: Histogram,
  }
  /**
   * Metric data for `exponential-histogram`.
   */
  export interface MetricDataExponentialHistogram {
    tag: 'exponential-histogram',
    val: ExponentialHistogram,
  }
  /**
   * Aggregated metrics data from an instrument.
   */
  export type AggregatedMetrics = AggregatedMetricsF64 | AggregatedMetricsU64 | AggregatedMetricsS64;
  /**
   * All metric data with `f64` value type.
   */
  export interface AggregatedMetricsF64 {
    tag: 'f64',
    val: MetricData,
  }
  /**
   * All metric data with `u64` value type.
   */
  export interface AggregatedMetricsU64 {
    tag: 'u64',
    val: MetricData,
  }
  /**
   * All metric data with `s64` value type.
   */
  export interface AggregatedMetricsS64 {
    tag: 's64',
    val: MetricData,
  }
  /**
   * `metric` is a collection of one or more aggregated time series from an instrument
   */
  export interface Metric {
    /**
     * The name of the instrument that created this data.
     */
    name: string,
    /**
     * The description of the instrument, which can be used in documentation.
     */
    description: string,
    /**
     * The unit in which the instrument reports.
     */
    unit: string,
    /**
     * The aggregated data from an instrument.
     */
    data: AggregatedMetrics,
  }
  /**
   * `scope-metrics` is a collection of `metric`s produced by a meter.
   */
  export interface ScopeMetrics {
    /**
     * The `instrumentation-scope` that the meter was created with.
     */
    scope: InstrumentationScope,
    /**
     * The list of aggregations created by the meter.
     */
    metrics: Array<Metric>,
  }
  /**
   * `resource-metrics` is a collection of `scope-metrics` and the associated `resource`
   * that created them.
   * 
   * See https://github.com/open-telemetry/opentelemetry-rust/blob/c811cde1ae21c624870c1b952190e687b16f76b8/opentelemetry-sdk/src/metrics/data/mod.rs#L13
   */
  export interface ResourceMetrics {
    /**
     * The entity that collected the metrics.
     */
    resource: Resource,
    /**
     * The collection of metrics with unique `instrumentation-scope`s.
     */
    scopeMetrics: Array<ScopeMetrics>,
  }
  /**
   * The WASI representation of the `OTelSdkError`.
   * 
   * See https://github.com/open-telemetry/opentelemetry-rust/blob/353bbb0d80fc35a26a00b4f4fed0dcaed23e5523/opentelemetry-sdk/src/error.rs#L15
   */
  export type OtelError = OtelErrorAlreadyShutdown | OtelErrorTimeout | OtelErrorInternalFailure;
  export interface OtelErrorAlreadyShutdown {
    tag: 'already-shutdown',
  }
  export interface OtelErrorTimeout {
    tag: 'timeout',
    val: Duration,
  }
  export interface OtelErrorInternalFailure {
    tag: 'internal-failure',
    val: string,
  }
}

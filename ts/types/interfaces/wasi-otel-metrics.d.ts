/// <reference path="./wasi-clocks-monotonic-clock.d.ts" />
/// <reference path="./wasi-clocks-wall-clock.d.ts" />
/// <reference path="./wasi-otel-tracing.d.ts" />
/// <reference path="./wasi-otel-types.d.ts" />
declare module 'wasi:otel/metrics@0.2.0-rc.2' {
  /**
   * Exports a resource's metric data.
   */
  export { _export as export };
  function _export(metrics: ResourceMetrics): void;
  export type Datetime = import('wasi:clocks/wall-clock@0.2.0').Datetime;
  export type Duration = import('wasi:clocks/monotonic-clock@0.2.0').Duration;
  export type KeyValue = import('wasi:otel/types@0.2.0-rc.2').KeyValue;
  export type InstrumentationScope = import('wasi:otel/types@0.2.0-rc.2').InstrumentationScope;
  export type Resource = import('wasi:otel/types@0.2.0-rc.2').Resource;
  export type SpanId = import('wasi:otel/tracing@0.2.0-rc.2').SpanId;
  export type TraceId = import('wasi:otel/tracing@0.2.0-rc.2').TraceId;
  /**
   * An error resulting from `export` being called.
   */
  export type Error = string;
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
   * This is the default temporality.
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
   * The number types available for any given instrument.
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
   * A single data point in a time series to be associated with a `gauge`.
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
   * A single data point in a time series to be associated with a `sum`.
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
   * A single data point in a time series to be associated with a `histogram`.
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
   * A single data point in a time series to be associated with an `exponential-histogram `.
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
  export type MetricData = MetricDataF64Gauge | MetricDataF64Sum | MetricDataF64Histogram | MetricDataF64ExponentialHistogram | MetricDataU64Gauge | MetricDataU64Sum | MetricDataU64Histogram | MetricDataU64ExponentialHistogram | MetricDataS64Gauge | MetricDataS64Sum | MetricDataS64Histogram | MetricDataS64ExponentialHistogram;
  /**
   * Metric data for an f64 gauge.
   */
  export interface MetricDataF64Gauge {
    tag: 'f64-gauge',
    val: Gauge,
  }
  /**
   * Metric data for an f64 sum.
   */
  export interface MetricDataF64Sum {
    tag: 'f64-sum',
    val: Sum,
  }
  /**
   * Metric data for an f64 histogram.
   */
  export interface MetricDataF64Histogram {
    tag: 'f64-histogram',
    val: Histogram,
  }
  /**
   * Metric data for an f64 exponential-histogram.
   */
  export interface MetricDataF64ExponentialHistogram {
    tag: 'f64-exponential-histogram',
    val: ExponentialHistogram,
  }
  /**
   * Metric data for an u64 gauge.
   */
  export interface MetricDataU64Gauge {
    tag: 'u64-gauge',
    val: Gauge,
  }
  /**
   * Metric data for an u64 sum.
   */
  export interface MetricDataU64Sum {
    tag: 'u64-sum',
    val: Sum,
  }
  /**
   * Metric data for an u64 histogram.
   */
  export interface MetricDataU64Histogram {
    tag: 'u64-histogram',
    val: Histogram,
  }
  /**
   * Metric data for an u64 exponential-histogram.
   */
  export interface MetricDataU64ExponentialHistogram {
    tag: 'u64-exponential-histogram',
    val: ExponentialHistogram,
  }
  /**
   * Metric data for an s64 gauge.
   */
  export interface MetricDataS64Gauge {
    tag: 's64-gauge',
    val: Gauge,
  }
  /**
   * Metric data for an s64 sum.
   */
  export interface MetricDataS64Sum {
    tag: 's64-sum',
    val: Sum,
  }
  /**
   * Metric data for an s64 histogram.
   */
  export interface MetricDataS64Histogram {
    tag: 's64-histogram',
    val: Histogram,
  }
  /**
   * Metric data for an s64 exponential-histogram.
   */
  export interface MetricDataS64ExponentialHistogram {
    tag: 's64-exponential-histogram',
    val: ExponentialHistogram,
  }
  /**
   * A collection of one or more aggregated time series from a metric.
   */
  export interface Metric {
    /**
     * The name of the metric that created this data.
     */
    name: string,
    /**
     * The description of the metric, which can be used in documentation.
     */
    description: string,
    /**
     * The unit in which the metric reports.
     */
    unit: string,
    /**
     * The aggregated data from a metric.
     */
    data: MetricData,
  }
  /**
   * A collection of `metric`s produced by a meter.
   */
  export interface ScopeMetrics {
    /**
     * The instrumentation scope that the meter was created with.
     */
    scope: InstrumentationScope,
    /**
     * The list of aggregations created by the meter.
     */
    metrics: Array<Metric>,
  }
  /**
   * A collection of `scope-metrics` and the associated `resource` that created them.
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
}

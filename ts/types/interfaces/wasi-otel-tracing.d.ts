declare module "wasi:otel/tracing@0.2.0-draft" {
  /**
   * Called when a span is started.
   */
  export function onStart(span: SpanData, parent: SpanContext): void;
  /**
   * Called when a span is ended.
   */
  export function onEnd(span: SpanData): void;
  /**
   * Returns the current span context of the host.
   */
  export function currentSpanContext(): SpanContext;

  export interface Datetime {
    seconds: bigint;
    nanoseconds: number;
  }
  /**
   * The trace that this `span-context` belongs to.
   *
   * 16 bytes encoded as a hexadecimal string.
   */
  export type TraceId = string;
  /**
   * The id of this `span-context`.
   *
   * 8 bytes encoded as a hexadecimal string.
   */
  export type SpanId = string;
  /**
   * Flags that can be set on a `span-context`.
   */
  export interface TraceFlags {
    /**
     * Whether the `span` should be sampled or not.
     */
    sampled?: boolean;
  }
  /**
   * Carries system-specific configuration data, represented as a list of key-value pairs. `trace-state` allows multiple tracing systems to participate in the same trace.
   *
   * If any invalid keys or values are provided then the `trace-state` will be treated as an empty list.
   */
  export type TraceState = Array<[string, string]>;
  /**
   * Identifying trace information about a span that can be serialized and propagated.
   */
  export interface SpanContext {
    /**
     * The `trace-id` for this `span-context`.
     */
    traceId: TraceId;
    /**
     * The `span-id` for this `span-context`.
     */
    spanId: SpanId;
    /**
     * The `trace-flags` for this `span-context`.
     */
    traceFlags: TraceFlags;
    /**
     * Whether this `span-context` was propagated from a remote parent.
     */
    isRemote: boolean;
    /**
     * The `trace-state` for this `span-context`.
     */
    traceState: TraceState;
  }
  /**
   * Describes the relationship between the Span, its parents, and its children in a trace.
   * # Variants
   *
   * ## `"client"`
   *
   * Indicates that the span describes a request to some remote service. This span is usually the parent of a remote server span and does not end until the response is received.
   * ## `"server"`
   *
   * Indicates that the span covers server-side handling of a synchronous RPC or other remote request. This span is often the child of a remote client span that was expected to wait for a response.
   * ## `"producer"`
   *
   * Indicates that the span describes the initiators of an asynchronous request. This parent span will often end before the corresponding child consumer span, possibly even before the child span starts. In messaging scenarios with batching, tracing individual messages requires a new producer span per message to be created.
   * ## `"consumer"`
   *
   * Indicates that the span describes a child of an asynchronous consumer request.
   * ## `"internal"`
   *
   * Default value. Indicates that the span represents an internal operation within an application, as opposed to an operations with remote parents or children.
   */
  export type SpanKind =
    | "client"
    | "server"
    | "producer"
    | "consumer"
    | "internal";
  /**
   * The key part of attribute `key-value` pairs.
   */
  export type Key = string;
  /**
   * The value part of attribute `key-value` pairs.
   */
  export type Value =
    | ValueString
    | ValueBool
    | ValueF64
    | ValueS64
    | ValueStringArray
    | ValueBoolArray
    | ValueF64Array
    | ValueS64Array;
  /**
   * A string value.
   */
  export interface ValueString {
    tag: "string";
    val: string;
  }
  /**
   * A boolean value.
   */
  export interface ValueBool {
    tag: "bool";
    val: boolean;
  }
  /**
   * A double precision floating point value.
   */
  export interface ValueF64 {
    tag: "f64";
    val: number;
  }
  /**
   * A signed 64 bit integer value.
   */
  export interface ValueS64 {
    tag: "s64";
    val: bigint;
  }
  /**
   * A homogeneous array of string values.
   */
  export interface ValueStringArray {
    tag: "string-array";
    val: Array<string>;
  }
  /**
   * A homogeneous array of boolean values.
   */
  export interface ValueBoolArray {
    tag: "bool-array";
    val: Array<boolean>;
  }
  /**
   * A homogeneous array of double precision floating point values.
   */
  export interface ValueF64Array {
    tag: "f64-array";
    val: Float64Array;
  }
  /**
   * A homogeneous array of 64 bit integer values.
   */
  export interface ValueS64Array {
    tag: "s64-array";
    val: BigInt64Array;
  }
  /**
   * A key-value pair describing an attribute.
   */
  export interface KeyValue {
    /**
     * The attribute name.
     */
    key: Key;
    /**
     * The attribute value.
     */
    value: Value;
  }
  /**
   * An event describing a specific moment in time on a span and associated attributes.
   */
  export interface Event {
    /**
     * Event name
     */
    name: string;
    /**
     * Event time
     */
    time: Datetime;
    /**
     * Event attributes
     */
    attributes: Array<KeyValue>;
  }
  /**
   * Describes a relationship to another `span`.
   */
  export interface Link {
    /**
     * Denotes which `span` to link to.
     */
    spanContext: SpanContext;
    /**
     * Attributes describing the link.
     */
    attributes: Array<KeyValue>;
  }
  /**
   * The `status` of a `span`.
   */
  export type Status = StatusUnset | StatusOk | StatusError;
  /**
   * The default status.
   */
  export interface StatusUnset {
    tag: "unset";
  }
  /**
   * The operation has been validated by an Application developer or Operator to have completed successfully.
   */
  export interface StatusOk {
    tag: "ok";
  }
  /**
   * The operation contains an error with a description.
   */
  export interface StatusError {
    tag: "error";
    val: string;
  }
  /**
   * Describes the instrumentation scope that produced a span.
   */
  export interface InstrumentationScope {
    /**
     * Name of the instrumentation scope.
     */
    name: string;
    /**
     * The library version.
     */
    version?: string;
    /**
     * Schema URL used by this library.
     * https://github.com/open-telemetry/opentelemetry-specification/blob/v1.9.0/specification/schemas/overview.md#schema-url
     */
    schemaUrl?: string;
    /**
     * Specifies the instrumentation scope attributes to associate with emitted telemetry.
     */
    attributes: Array<KeyValue>;
  }
  /**
   * The data associated with a span.
   */
  export interface SpanData {
    /**
     * Span context
     */
    spanContext: SpanContext;
    /**
     * Span parent id
     * TODO: No clue what this is for
     */
    parentSpanId: string;
    /**
     * Span kind
     */
    spanKind: SpanKind;
    /**
     * Span name
     */
    name: string;
    /**
     * Span start time
     */
    startTime: Datetime;
    /**
     * Span end time
     */
    endTime: Datetime;
    /**
     * Span attributes
     */
    attributes: Array<KeyValue>;
    /**
     * Span events
     */
    events: Array<Event>;
    /**
     * Span Links
     */
    links: Array<Link>;
    /**
     * Span status
     */
    status: Status;
    /**
     * Instrumentation scope that produced this span
     */
    instrumentationScope: InstrumentationScope;
  }
}

/// <reference path="./wasi-clocks-wall-clock.d.ts" />
/// <reference path="./wasi-otel-tracing.d.ts" />
/// <reference path="./wasi-otel-types.d.ts" />
declare module 'wasi:otel/logs@0.2.0-draft' {
  /**
   * Called when a log is emitted.
   */
  export function onEmit(data: LogRecord): void;
  export type InstrumentationScope = import('wasi:otel/types@0.2.0-draft').InstrumentationScope;
  export type Resource = import('wasi:otel/types@0.2.0-draft').Resource;
  export type Value = import('wasi:otel/types@0.2.0-draft').Value;
  export type KeyValue = import('wasi:otel/types@0.2.0-draft').KeyValue;
  export type SpanId = import('wasi:otel/tracing@0.2.0-draft').SpanId;
  export type TraceId = import('wasi:otel/tracing@0.2.0-draft').TraceId;
  export type TraceFlags = import('wasi:otel/tracing@0.2.0-draft').TraceFlags;
  export type Datetime = import('wasi:clocks/wall-clock@0.2.0').Datetime;
  /**
   * Represents the recording of an event.
   */
  export interface LogRecord {
    /**
     * Time when the event occurred.
     */
    timestamp?: Datetime,
    /**
     * Time when the event was observed.
     */
    observedTimestamp?: Datetime,
    /**
     * The severity text(also known as log level).
     */
    severityText?: string,
    /**
     * The numerical value of the severity ranging from 1-24.
     */
    severityNumber?: number,
    /**
     * The body of the log record.
     */
    body?: Value,
    /**
     * Additional information about the specific event occurrence.
     */
    attributes?: Array<KeyValue>,
    /**
     * Name that identifies the class / type of event.
     */
    eventName?: string,
    /**
     * Describes the source of the log.
     */
    resource?: Resource,
    /**
     * Describes the scope that emitted the log.
     */
    instrumentationScope?: InstrumentationScope,
    /**
     * Request trace id.
     */
    traceId?: TraceId,
    /**
     * Request span id.
     */
    spanId?: SpanId,
    /**
     * W3C trace flag.
     */
    traceFlags?: TraceFlags,
  }
}

import { KeyValue, Resource, InstrumentationScope } from './wasi-otel-types';

export type Datetime = { seconds: bigint; nanoseconds: number };
export type SpanId = string;
export type TraceId = string;
export type TraceFlags = number;
export type Value = string;

export interface LogRecord {
  timestamp?: Datetime;
  observedTimestamp?: Datetime;
  severityText?: string;
  severityNumber?: number;
  body?: Value;
  attributes?: Array<KeyValue>;
  eventName?: string;
  resource?: Resource;
  instrumentationScope?: InstrumentationScope;
  traceId?: TraceId;
  spanId?: SpanId;
  traceFlags?: TraceFlags;
}

export function onEmit(_data: LogRecord): void {
  // Mock implementation - does nothing in tests
}

export { KeyValue };

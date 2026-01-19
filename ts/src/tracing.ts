import {
  onStart as wasiSpanStart,
  onEnd as wasiSpanEnd,
  outerSpanContext as wasiOuterSpanContext,
  SpanContext as WasiSpanContext,
  TraceState as WasiTraceState,
  TraceFlags as WasiTraceFlags,
  SpanData as WasiSpanData,
  SpanKind as WasiSpanKind,
} from 'wasi:otel/tracing@0.2.0-draft';
import {
  ReadableSpan,
  Span,
  SpanProcessor,
} from '@opentelemetry/sdk-trace-base';
import {
  Context,
  trace,
  createTraceState,
  SpanContext,
  TraceFlags,
  TraceState,
} from '@opentelemetry/api';
import {
  dateTimeToWasi,
  attributesToWasi,
  instrumentationScopeToWasi,
} from './types';

export class WasiTraceContextPropagator {
  constructor() {}
  /**
   * Retrieves trace context from a WASI host and combines it with the current trace context.
   * @param cx The current trace context.
   * @returns The combined host and current trace context.
   */
  extract(cx: Context): Context {
    return trace.setSpanContext(cx, wasiToSpanContext(wasiOuterSpanContext()));
  }
}

export class WasiSpanProcessor implements SpanProcessor {
  async forceFlush(): Promise<void> {
    // no-op
  }

  onStart(span: Span, _: Context): void {
    wasiSpanStart(spanContextToWasi(span.spanContext()));
  }

  onEnd(span: ReadableSpan): void {
    wasiSpanEnd(spanDataToWasi(span));
  }

  async shutdown(): Promise<void> {
    // no-op
  }
}

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
  };
}

/**
 * Converts WASI TraceState to OpenTelemetry TraceState
 */
function wasiToTraceState(wts: WasiTraceState): TraceState {
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
  };
}

/**
 * Converts OpenTelemetry TraceFlags to WASI TraceFlags
 */
function traceFlagsToWasi(flags: TraceFlags): WasiTraceFlags {
  const SAMPLED = 0x01;
  return (flags & SAMPLED) === SAMPLED ? { sampled: true } : { sampled: false };
}

/**
 * Converts OpenTelemetry TraceState to WASI TraceState
 */
function traceStateToWasi(value: TraceState | undefined): WasiTraceState {
  if (value == undefined) {
    return [];
  }

  const serialized = value.serialize();
  if (serialized === '') {
    return [];
  }

  return serialized.split(',').map((entry): [string, string] => {
    // This ensures that a pattern like "foo=bar=baz" is split into
    // Key("foo"), Value("bar=baz")
    const [key, ...rest] = entry.split('=');
    const value = rest.join('=');
    return [key, value];
  });
}

/**
 * Converts OpenTelemetry ReadableSpan to WASI SpanData
 */
export function spanDataToWasi(span: ReadableSpan): WasiSpanData {
  return {
    name: span.name,
    startTime: dateTimeToWasi(span.startTime),
    spanContext: spanContextToWasi(span.spanContext()),
    parentSpanId: span.parentSpanContext?.spanId || '',
    spanKind: ['internal', 'server', 'client', 'producer', 'consumer'][
      span.kind
    ] as WasiSpanKind,
    endTime: dateTimeToWasi(span.endTime),
    attributes: attributesToWasi(span.attributes),
    events: span.events.map((e) => ({
      name: e.name,
      time: dateTimeToWasi(e.time),
      attributes: attributesToWasi(e.attributes),
    })),
    links: span.links.map((link) => ({
      spanContext: spanContextToWasi(link.context),
      attributes: attributesToWasi(link.attributes),
    })),
    status:
      span.status.code === 0
        ? { tag: 'unset' }
        : span.status.code === 1
          ? { tag: 'ok' }
          : { tag: 'error', val: span.status.message || '' },
    instrumentationScope: instrumentationScopeToWasi(span.instrumentationScope),
    droppedAttributes: span.droppedAttributesCount,
    droppedEvents: span.droppedEventsCount,
    droppedLinks: span.droppedLinksCount,
  };
}

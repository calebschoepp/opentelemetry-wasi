import {
  TraceFlags as WasiTraceFlags,
  TraceState as WasiTraceState,
  SpanContext as WasiSpanContext,
  Datetime,
} from "wasi:otel/tracing@0.2.0-draft";
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
    TraceState
} from "@opentelemetry/api";

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
 * Converts OpenTelemetry TraceState to WASI TraceTraceState
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

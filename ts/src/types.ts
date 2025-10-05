import {
  TraceFlags as wasiTraceFlags,
  TraceState as wasiTraceState,
  SpanContext as wasiSpanContext,
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

export function spanContextToWasi(ctx: SpanContext): wasiSpanContext {
  return {
    traceId: ctx.traceId,
    spanId: ctx.spanId,
    traceFlags: traceFlagsToWasi(ctx.traceFlags),
    isRemote: ctx.isRemote || false,
    traceState: traceStateToWasi(ctx.traceState),
  }
}

export function wasiToSpanContext(ctx: wasiSpanContext): SpanContext {
    const wasiToTraceState = (wts: wasiTraceState): TraceState => {
        let traceState = createTraceState();
        for (const [key, value] of wts) {
            traceState = traceState.set(key, value);
        }

        return traceState;
    };

    return {
        traceId: ctx.traceId,
        spanId: ctx.spanId,
        traceFlags: ctx.traceFlags.sampled ? TraceFlags.SAMPLED : TraceFlags.NONE,
        isRemote: ctx.isRemote || false,
        traceState: wasiToTraceState(ctx.traceState),
    }
}

export function attributesToWasi(attrs: Attributes | undefined): KeyValue[] {
  if (attrs === undefined) {
    return [];
  }
  return Object.entries(attrs).map(([key, value]) => ({
    key,
    value: attributeValueToWasi(value),
  }));
}

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

function traceFlagsToWasi(flags: TraceFlags): wasiTraceFlags {
  const SAMPLED = 0x01;
  return (flags & SAMPLED) === SAMPLED ? {sampled: true} : {sampled: false};
}

function traceStateToWasi(value: TraceState | undefined): wasiTraceState {
  if (value == undefined) {
    return [];
  }
  return value.serialize().split(',').map(entry => {
     // TODO: I'm attempting to mimic rust's `split_once` method
      const [key, ...rest] = entry.split('=');
      const value = rest.join('=');
      return [key, value];
  });
}

export function dateTimeToWasi(time: HrTime): Datetime {
  return {
    seconds: BigInt(time[0]),
    nanoseconds: time[1],
  };
}
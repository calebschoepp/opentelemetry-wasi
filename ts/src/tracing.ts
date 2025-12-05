import {
  onStart as wasiSpanStart,
  onEnd as wasiSpanEnd,
  outerSpanContext as wasiOuterSpanContext,
} from "wasi:otel/tracing@0.2.0-draft";
import {
  ReadableSpan,
  Span,
  SpanProcessor,
} from "@opentelemetry/sdk-trace-base";
import {
    Context,
    trace,
} from "@opentelemetry/api";
import {
    spanContextToWasi,
    spanDataToWasi,
    wasiToSpanContext,
} from "./types"

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

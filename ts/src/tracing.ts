import {
  onStart as wasiSpanStart,
  onEnd as wasiSpanEnd,
  SpanKind,
  outerSpanContext,
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
    attributesToWasi,
    dateTimeToWasi,
    wasiToSpanContext,
} from "./types"

export class WasiTraceContextPropagator {
    constructor() {}
    extract(cx: Context): Context {
        return trace.setSpanContext(cx, wasiToSpanContext(outerSpanContext()));
    }
}

export class WasiSpanProcessor implements SpanProcessor {
  forceFlush(): Promise<void> {
    throw new Error("Method not implemented.");
  }

  onStart(span: Span, _: Context): void {
    wasiSpanStart(spanContextToWasi(span.spanContext()));
  }

  onEnd(span: ReadableSpan): void {
    // TODO: I'm unclear on how to reparent the guest spans...
    let ctx = span.spanContext();
    console.log("SpanContext:", ctx)
    wasiSpanEnd({
      name: span.name,
      startTime: dateTimeToWasi(span.startTime),
      spanContext: spanContextToWasi(ctx),
      parentSpanId: span.parentSpanId || "",
      spanKind: ["internal", "server", "client", "producer", "consumer"][span.kind] as SpanKind,
      endTime:dateTimeToWasi(span.endTime),
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
        name: span.instrumentationLibrary.name,
        version: span.instrumentationLibrary.version,
        schemaUrl: span.instrumentationLibrary.schemaUrl,
        // Although other SDKs use the InstrumentationScope.attributes field;
        // the `opentelemetry-js` SDK does not.
        // See https://github.com/open-telemetry/opentelemetry-js/blob/06621d27068881cc45329ecc76564f1d0c0b133f/packages/opentelemetry-core/src/common/types.ts#L47
        attributes: [],
      },
      droppedAttributes: span.droppedAttributesCount,
      droppedEvents: span.droppedEventsCount,
      droppedLinks: span.droppedLinksCount,
    });
  }

  shutdown(): Promise<void> {
    // Do not care
    throw new Error("Method not implemented.");
  }
}
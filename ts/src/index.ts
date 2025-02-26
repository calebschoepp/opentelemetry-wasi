import {
  onStart as wasiSpanStart,
  onEnd as wasiSpanEnd,
  //@ts-ignore
} from "wasi:otel/tracing@0.2.0-draft";
import {
  ReadableSpan,
  Span,
  SpanProcessor,
} from "@opentelemetry/sdk-trace-base";
import { Context } from "@opentelemetry/api";

export class WasiProcessor implements SpanProcessor {
  forceFlush(): Promise<void> {
    throw new Error("Method not implemented.");
  }

  onStart(span: Span, parentContext: Context): void {
    wasiSpanStart(
      {
        name: "foo",
        startTime: {
          seconds: BigInt(span.startTime[0]),
          nanoseconds: span.startTime[1],
        },
        spanContext: {
          traceId: "",
          spanId: "",
          traceFlags: { sampled: true },
          isRemote: false,
          traceState: [],
        },
        parentSpanId: "",
        spanKind: "client",
        endTime: {
          seconds: BigInt(span.endTime[0]),
          nanoseconds: span.endTime[1],
        },
        attributes: [],
        events: [],
        links: [],
        status: { tag: "unset" },
        instrumentationScope: {
          name: "",
          attributes: [],
        },
      },
      {
        traceId: "",
        spanId: "",
        traceFlags: { sampled: true },
        isRemote: false,
        traceState: [],
      }
    );
  }

  onEnd(span: ReadableSpan): void {
    wasiSpanEnd({
      name: "foo",
      startTime: {
        seconds: BigInt(span.startTime[0]),
        nanoseconds: span.startTime[1],
      },
      spanContext: {
        traceId: "",
        spanId: "",
        traceFlags: { sampled: true },
        isRemote: false,
        traceState: [],
      },
      parentSpanId: "",
      spanKind: "client",
      endTime: {
        seconds: BigInt(span.endTime[0]),
        nanoseconds: span.endTime[1],
      },
      attributes: [],
      events: [],
      links: [],
      status: { tag: "unset" },
      instrumentationScope: {
        name: "",
        attributes: [],
      },
    });
  }

  shutdown(): Promise<void> {
    // Do not care
    throw new Error("Method not implemented.");
  }
}

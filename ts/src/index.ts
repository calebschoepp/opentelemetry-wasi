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
import { BindOnceFuture } from "@opentelemetry/core";

export class WasiProcessor implements SpanProcessor {
  private _shutdownOnce: BindOnceFuture<void>;

  constructor() {
    this._shutdownOnce = new BindOnceFuture(this._shutdown, this);
  }

  forceFlush(): Promise<void> {
    // no-op
    return Promise.resolve();
  }

  onStart(span: Span, _: Context): void {
    if (this._shutdownOnce.isCalled) {
      return;
    }

    // TODO
    wasiSpanStart(
      // {
      //   name: "foo",
      //   startTime: {
      //     seconds: BigInt(span.startTime[0]),
      //     nanoseconds: span.startTime[1],
      //   },
      //   spanContext: {
      //     traceId: "",
      //     spanId: "",
      //     traceFlags: { sampled: true },
      //     isRemote: false,
      //     traceState: [],
      //   },
      //   parentSpanId: "",
      //   spanKind: "client",
      //   endTime: {
      //     seconds: BigInt(span.endTime[0]),
      //     nanoseconds: span.endTime[1],
      //   },
      //   attributes: [],
      //   events: [],
      //   links: [],
      //   status: { tag: "unset" },
      //   instrumentationScope: {
      //     name: "",
      //     attributes: [],
      //   },
      // },
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
    if (this._shutdownOnce.isCalled) {
      return;
    }

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
    return this._shutdownOnce.call();
  }

  private _shutdown(): void {
    // no-op
  }
}

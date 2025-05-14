/// <reference path="./interfaces/wasi-clocks-wall-clock.d.ts" />
/// <reference path="./interfaces/wasi-otel-tracing.d.ts" />
declare module 'wasi:otel/imports@0.2.0-draft' {
  export type * as WasiClocksWallClock020 from 'wasi:clocks/wall-clock@0.2.0'; // import wasi:clocks/wall-clock@0.2.0
  export type * as WasiOtelTracing020Draft from 'wasi:otel/tracing@0.2.0-draft'; // import wasi:otel/tracing@0.2.0-draft
}

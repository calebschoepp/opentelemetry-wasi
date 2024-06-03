# Processor

Example applications using the processor-like WIT interface.

## Pros

- This WIT interface can support both the OTel processor and tracer interfaces.

## Cons

- Doesn't work with Rust tracing which uses the tracer at the end of the span.

## Interface

```wit
interface observe {
  use wasi:clocks/wall-clock@0.2.0.{datetime};

  // TODO
  on-span-start: func(span-context: span-context);

  // TODO
  on-span-end: func(span: read-only-span);

  // TODO: Document.
  record read-only-span {
    // Span name.
    name: string,

    // Span context.
    span-context: span-context,

    // Span parent id.
    parent-span-id: string,

    // Span kind.
    span-kind: span-kind,

    // Span start time.
    start-time: datetime,

    // Span end time.
    end-time: datetime,

    // Span attributes. TODO: Support multiple types
    attributes: list<tuple<string, string>>,

    // Span resource.
    otel-resource: otel-resource,

    // TODO: Support dropped_attributes_count, events, links, status, and instrumentation lib
  }

  // Identifying trace information about a span.
  record span-context {
    // Hexidecimal representation of the trace id.
    trace-id: string,

    // Hexidecimal representation of the span id.
    span-id: string,

    // Hexidecimal representation of the trace flags
    trace-flags: string,

    // Span remoteness
    is-remote: bool,

    // Entirity of tracestate
    trace-state: string,
  }

  // TODO: Document this and children.
  enum span-kind {
    client,
    server,
    producer,
    consumer,
    internal
  }

  // An immutable representation of the entity producing telemetry as attributes.
  record otel-resource {
    // Resource attributes.
    attrs: list<tuple<string, string>>,

    // Resource schema url.
    schema-url: option<string>,
  }
}
```

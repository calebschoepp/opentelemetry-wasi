# Exporter

Example applications using the exporter-like WIT interface.

## Pros

- Simple WIT interface and host implementation.

## Cons

- Can't properly set parent span of guest -> host spans.

## Interface

```wit
interface observe {
  use wasi:clocks/wall-clock@0.2.0.{datetime};

  // Emit a given completed read-only-span to the o11y host.
  emit-span: func(span: read-only-span) -> result<_, string>;

  // get-parent-span-context returns the parent span context of the host.
  get-parent-span-context: func() -> span-context;

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

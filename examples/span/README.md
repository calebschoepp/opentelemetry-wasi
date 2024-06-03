# Span

Example applications using the span-like WIT interface.

## Pros

- Resource semantics work well for directly using in components.

## Cons

- Doesn't work with Rust tracing which uses the tracer at the end of the span.
- Requires re-writing the most of OTel.

## Interface

```wit
interface observe {
  use wasi:clocks/wall-clock@0.2.0.{datetime};

  // TODO: Document.
  resource span {
    // enter returns a new span with the given name.
    enter: static func(name: string) -> span;

    // set-attribute sets an attribute on the span.
    set-attribute: func(key: string, value: string);

    // close closes the span.
    close: func();
  }
```

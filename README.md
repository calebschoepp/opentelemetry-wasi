# OpenTelemetry WASI

Libraries to enable using OpenTelemetry within WebAssembly components backed by [WASI OTel](https://github.com/calebschoepp/wasi-otel).

## Concepts

To make it easier for folks without OpenTelemetry experience to follow along, we created a high-level overview of the different parts of OpenTelemetry relevant to this project. If you want more specifics the [OpenTelemetry Specs](https://opentelemetry.io/docs/specs/otel/overview/) is a great resource. 

### Tracing

[OpenTelemetry Tracing API Specs](https://opentelemetry.io/docs/specs/otel/trace/api/)

### Metrics

[OpenTelemetry Metrics API Specs](https://opentelemetry.io/docs/specs/otel/metrics/api/)

#### Terms
- `MeterProvider`: the entry point of the API. It provides access to `Meters`.
- `Meter`: responsible for creating `Instruments`.
- `Instrument`: responsible for reporting `Measurements`. An instrument defines what a measurement means/how it is used. There are a few different ways to report measurements: counters, histograms, gauges, and their synchronous/asynchronous variants.
- `Measurement`: a data point to be reported.

In summary, a measurement is collected by an instrument, which is organized under a meter, which is organized under a default or a custom MeterProvider.

### Logs

[OpenTelemetry Logs API Specs](https://opentelemetry.io/docs/specs/otel/logs/api/)
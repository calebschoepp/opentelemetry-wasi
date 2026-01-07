package main

import (
	"context"
	"fmt"
	"net/http"

	wasiLogs "github.com/calebschoepp/opentelemetry-wasi/logs"
	wasiMetrics "github.com/calebschoepp/opentelemetry-wasi/metrics"
	wasiTracing "github.com/calebschoepp/opentelemetry-wasi/tracing"
	spinhttp "github.com/spinframework/spin-go-sdk/v3/http"
	"github.com/spinframework/spin-go-sdk/v3/kv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	logApi "go.opentelemetry.io/otel/log"
	metricApi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func init() {
	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		/*
			### LOGS ###
		*/

		loggerProvider := log.NewLoggerProvider(log.WithProcessor(wasiLogs.NewWasiLogProcessor()))
		logger := loggerProvider.Logger("spin-logs")
		logRecord := logApi.Record{}
		logRecord.SetBody(logApi.StringValue("Hello from Go!"))
		logRecord.SetSeverity(logApi.SeverityInfo)
		logger.Emit(ctx, logRecord)

		/*
			### METRICS ###
		*/
		exporter := wasiMetrics.NewWasiMetricExporter()
		defer exporter.Export(ctx) // Export metrics to the host

		meterProvider := metric.NewMeterProvider(metric.WithReader(exporter))
		meter := meterProvider.Meter("spin-metrics")

		attrs := metricApi.WithAttributes(
			attribute.Key("spinkey1").String("spinvalue1"),
			attribute.Key("spinkey2").String("spinvalue2"),
		)

		counter, err := meter.Int64Counter("spin-counter")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		counter.Add(ctx, 10, attrs)

		upDownCounter, err := meter.Int64UpDownCounter("spin-up-down-counter")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		upDownCounter.Add(ctx, -1, attrs)

		histogram, err := meter.Int64Histogram("spin-histogram")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		histogram.Record(ctx, 9, attrs)
		histogram.Record(ctx, 15, attrs)

		gauge, err := meter.Float64Gauge("spin-gauge")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		gauge.Record(ctx, 123.456, attrs)

		/*
			### TRACING ###
		*/

		tracerProvider := trace.NewTracerProvider(trace.WithSpanProcessor(wasiTracing.NewWasiSpanProcessor()))
		propagator := wasiTracing.NewTraceContextPropagator()
		hostCtx := propagator.Extract(ctx)
		otel.SetTracerProvider(tracerProvider)

		tracer := tracerProvider.Tracer("spin-tracer")

		func() {
			mainCtx, mainSpan := tracer.Start(hostCtx, "main-operation")
			defer mainSpan.End()

			_, childSpan := tracer.Start(mainCtx, "child-operation")
			defer childSpan.End()

			store, err := kv.OpenDefault()
			if err != nil {
				childSpan.RecordError(err)
				childSpan.SetStatus(codes.Error, "failed to open kv store")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := store.Set("foo", []byte("bar")); err != nil {
				childSpan.RecordError(err)
				childSpan.SetStatus(codes.Error, "failed to set value")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			childSpan.SetStatus(codes.Ok, "success")
		}()

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "Hello World!")
	})
}

func main() {}

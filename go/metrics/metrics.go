package metrics

import (
	"context"

	wasiMetrics "github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// A metric exporter that sends OpenTelemetry metrics to a WASI host.
type WasiMetricExporter struct {
	metric.Reader
}

// Create a new WasiMetricExporter
func NewWasiMetricExporter() *WasiMetricExporter {
	return &WasiMetricExporter{
		Reader: metric.NewManualReader(),
	}
}

// Exports metric data to a compatible host or component.
func (e *WasiMetricExporter) Export(ctx context.Context) {
	metrics := metricdata.ResourceMetrics{}
	if err := e.Collect(ctx, &metrics); err != nil {
		// This gives the user control over how internal errors are handled
		otel.Handle(err)
		return
	}

	wasiMetrics.Export(toWasiResourceMetrics(metrics))
}

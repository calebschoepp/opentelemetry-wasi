package logs

import (
	"context"

	wasiLogs "github.com/calebschoepp/opentelemetry-wasi/internal/wasi_otel_logs"
	"go.opentelemetry.io/otel/sdk/log"
)

type WasiLogProcessor struct {
	log.Processor
}

func NewWasiLogProcessor() WasiLogProcessor {
	return WasiLogProcessor{
		Processor: newWasiLogExporter(),
	}
}

type wasiLogExporter struct {
	log.Exporter
}

func newWasiLogExporter() *wasiLogExporter {
	return &wasiLogExporter{}
}

func (w *wasiLogExporter) Enabled(ctx context.Context, params log.EnabledParameters) bool {
	return true
}

func (w *wasiLogExporter) OnEmit(ctx context.Context, record *log.Record) error {
	wasiLogs.OnEmit(toWasiLogRecord(*record))
	return nil
}

func (w *wasiLogExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (w *wasiLogExporter) ForceFlush(ctx context.Context) error {
	return nil
}

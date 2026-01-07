package logs

import (
	"context"

	wasiLogs "github.com/calebschoepp/opentelemetry-wasi/wit_component/wasi_otel_logs"
	"go.opentelemetry.io/otel/sdk/log"
)

type WasiLogProcessor struct {
	log.Processor
}

func NewWasiLogProcessor() WasiLogProcessor {
	return WasiLogProcessor{
		Processor: log.NewSimpleProcessor(newWasiLogExporter()),
	}
}

type wasiLogExporter struct {
	log.Exporter
}

func newWasiLogExporter() *wasiLogExporter {
	return &wasiLogExporter{}
}

func (w *wasiLogExporter) Export(ctx context.Context, records []log.Record) error {
	for _, record := range records {
		wasiLogs.OnEmit(toWasiLogRecord(record))
	}

	return nil
}

func (w *wasiLogExporter) Shutdown(ctx context.Context) error {
	return nil
}

func (w *wasiLogExporter) ForceFlush(ctx context.Context) error {
	return nil
}

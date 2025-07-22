package main

import (
	"fmt"
	"net/http"

	"github.com/ydnar/wasi-http-go/wasihttp"

	otelWasi "github.com/calebschoepp/opentelemetry-wasi"
	spinkv "github.com/spinframework/spin-go-sdk/v2/kv"
	"github.com/spinframework/spin-go-sdk/v2/wit"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

var _ = wit.Wit

func init() {
	wasihttp.Serve(&WasiHTTP{})
}

type WasiHTTP struct{}

func (WasiHTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wasiProcessor := otelWasi.NewWasiProcessor()
	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithSpanProcessor(wasiProcessor),
	)

	otel.SetTracerProvider(tracerProvider)
	tracer := otel.Tracer("basic-spin")

	wasiPropagator := otelWasi.NewTraceContextPropagator()
	wasiPropagator.Extract(r.Context())

	ctx, span := tracer.Start(r.Context(), "main-operation")
	defer span.End()

	span.SetAttributes(attribute.String("my-attribute", "my-value"))
	span.AddEvent(
		"Main span event",
		trace.WithAttributes(attribute.String("foo", "1")),
	)

	_, childSpan := tracer.Start(ctx, "child-operation")
	childSpan.AddEvent(
		"Sub span event",
		trace.WithAttributes(attribute.String("bar", "1")),
	)
	defer childSpan.End()

	store, err := spinkv.OpenDefault()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to open kv store"))
	}

	if err := store.Set("foo", []byte("bar")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to write to kv store"))
	}

	w.Header().Set("content-type", "text/plain")
	fmt.Fprintln(w, "Hello, world!")
}

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	spinhttp "github.com/fermyon/spin/sdk/go/v2/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const name = "go.opentelemetry.io/otel/example/dice"

var (
	otelInitOnce sync.Once
	tracer       = otel.Tracer(name)
)

func init() {
	fmt.Println("0")
	otelInitOnce.Do(func() {
		if _, err := setupOtelSDK(context.Background()); err != nil {
			fmt.Printf("Failed to setup OpenTelemetry: %v\n", err)
		}

		tracer = otel.Tracer("spin-basic")
	})

	fmt.Println("1")

	spinhttp.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/plain")
		fmt.Fprintln(w, "Hello, world!")
		fmt.Fprintln(w, "2")
		ctx := r.Context()
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(r.Header))
		fmt.Fprintln(w, "3")
		ctx, span := tracer.Start(
			ctx,
			"HTTP"+r.Method,
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.route", "/"),
			),
		)
		defer span.End()

		fmt.Fprintln(w, "4")

		r = r.WithContext(ctx)

		fmt.Fprintln(w, "5")
		diceHandler(w, r)
		fmt.Fprintln(w, "6")

		span.SetAttributes(attribute.Int("http.status_code", 200))
		fmt.Fprintln(w, "7")
	})
	fmt.Println("Request handled")
}

func rollDice(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "roll")
	defer span.End()

	roll := 1 + rand.Intn(6)

	var msg string
	if player := r.PathValue("player"); player != "" {
		msg = fmt.Sprintf("%s is rolling the dice", player)
	} else {
		msg = "Anonymous player is rolling the dice"
	}

	fmt.Fprintln(w, msg)

	rollValueAttr := attribute.Int("roll.value", roll)
	span.SetAttributes(rollValueAttr)

	if _, err := io.WriteString(w, strconv.Itoa(roll)+"\n"); err != nil {
		fmt.Fprintf(w, "Write failed: %v\n", err)
	}
}

func diceHandler(w http.ResponseWriter, r *http.Request) {
	// You can also manually create spans for more detailed tracing
	ctx := r.Context()
	_, span := tracer.Start(ctx, "roll_dice_operation")
	defer span.End()

	w.Header().Set("Content-Type", "text/plain")
	rollDice(w, r)
}

func setupOtelSDK(ctx context.Context) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	prop := newPropogator()
	otel.SetTextMapPropagator(prop)

	tracerProvider, err := newTracerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	return
}

func newPropogator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider() (*sdkTrace.TracerProvider, error) {
	traceExporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tracerProvider := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(
			traceExporter,
			sdkTrace.WithBatchTimeout(time.Second),
		),
	)

	return tracerProvider, nil
}

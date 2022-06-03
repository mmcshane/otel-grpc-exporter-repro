package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	otel.SetErrorHandler(
		otel.ErrorHandlerFunc(func(err error) {
			log.Fatalf("failed with %q", err) // kill process on error
		}),
	)
	exp, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("localhost:4317"),

		// unrelated - just speeds up the repro
		otlptracegrpc.WithTimeout(1*time.Second),

		// this _should_ be enough to connect to a non-TLS server but isn't
		otlptracegrpc.WithDialOption(grpc.WithInsecure()),

		// Uncommenting this to allow export to non-TLS server
		// otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(sdktrace.NewSimpleSpanProcessor(exp)),
	)

	ctx, span := tp.Tracer("main").Start(ctx, "a_span")
	time.Sleep(100 * time.Millisecond)
	span.End()

	tp.Shutdown(ctx)
	fmt.Println("succeeded")
}

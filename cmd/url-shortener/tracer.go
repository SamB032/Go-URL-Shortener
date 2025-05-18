package main

import (
	"context"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func initTracer(otelExporterEndpoint string) (*sdktrace.TracerProvider, error) {
	exp, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(strings.TrimPrefix(otelExporterEndpoint, "http://")),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-url-shortener"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}

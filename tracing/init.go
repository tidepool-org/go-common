package tracing

import (
	"log"

	"go.opentelemetry.io/otel"
	otelglobal "go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/propagators"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func InitTracer() {
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}
	cfg := sdktrace.Config{
		DefaultSampler: sdktrace.AlwaysSample(),
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(cfg),
		sdktrace.WithSyncer(exporter),
	)
	if err != nil {
		log.Fatal(err)
	}
	otelglobal.SetTracerProvider(tp)
	otelglobal.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(propagators.TraceContext{}))
}

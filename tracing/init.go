package tracing

// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Example using the OTLP exporter + collector + third-party backends. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/stdout"

	"go.opentelemetry.io/otel/propagators"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"

	sdkmetric "go.opentelemetry.io/otel/sdk/export/metric"
)

//TraceConfig defines constants for configuring the trace system
type TraceConfig struct {
	Collector    string `envconfig:"OTEL_COLLECTOR_HOST" required:"true"`
	PodName      string `envconfig:"POD_NAME" required:"true"`
	PodNamespace string `envconfig:"POD_NAMESPACE" required:"true"`
	PodIP        string `envconfig:"POD_IP" required:"true"`
}

//SamplingProvider provides the sampling algorithm for the tracer
func SamplingProvider() sdktrace.Sampler {
	return sdktrace.AlwaysSample()
}

//TraceConfigProvider provides the location of the tracing collector host and info about the current process
func TraceConfigProvider() (traceConfig TraceConfig, err error) {
	if err = envconfig.Process("", &traceConfig); err != nil {
		return TraceConfig{}, err
	}
	return traceConfig, nil
}

//ExporterProvider provides an oltp exporter
func ExporterProvider(traceConfig TraceConfig) (*otlp.Exporter, error) {
	exp, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress(traceConfig.Collector),
		otlp.WithGRPCDialOption(grpc.WithBlock()))

	return exp, err
}

// TraceExporter converts an otlp Exporter into a metric exporter
func TraceExporter(exp *otlp.Exporter) sdkmetric.Exporter {
	return exp
}

// MetricExporter converts an otlp Exporter into a span exporter
func MetricExporter(exp *otlp.Exporter) export.SpanExporter {
	return exp
}

// StdoutMetricExporter converts an otlp Exporter into a span exporter
func StdoutMetricExporter(exp *stdout.Exporter) export.SpanExporter {
	return exp
}

// StdoutTraceExporter converts an otlp Exporter into a metric exporter
func StdoutTraceExporter(exp *stdout.Exporter) sdkmetric.Exporter {
	return exp
}

//StdoutExporterProvider provides an exporter that writes to the stdout
func StdoutExporterProvider(traceConfig TraceConfig) (*stdout.Exporter, error) {
	exporter, err := stdout.NewExporter(stdout.WithPrettyPrint())
	return exporter, err
}

//SpanProcessorProvider provides a span processor
func SpanProcessorProvider(exp export.SpanExporter) sdktrace.SpanProcessor {
	return sdktrace.NewBatchSpanProcessor(exp)
}

//TracerProvider provides a tracer
func TracerProvider(traceConfig TraceConfig, bsp sdktrace.SpanProcessor, sampler sdktrace.Sampler) trace.TracerProvider {
	res := resource.New(
		semconv.ServiceNameKey.String(traceConfig.PodName),
		semconv.K8SNamespaceNameKey.String(traceConfig.PodNamespace),
		semconv.ServiceInstanceIDKey.String(traceConfig.PodIP),
	)

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	return provider
}

//PushProvider provides an OpenTelemetry push controller
func PushProvider(exp sdkmetric.Exporter) *push.Controller {
	return push.New(
		basic.New(simple.NewWithExactDistribution(), exp),
		exp,
		push.WithPeriod(2*time.Second),
	)
}

//TextMapPropagatorProvider provides a text map propagator
func TextMapPropagatorProvider() otel.TextMapPropagator {
	return otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{})
}

//MetricProvider provides a meter provider
func MetricProvider(pusher *push.Controller) metric.MeterProvider {
	return pusher.MeterProvider()
}

// Params needed to start Tracer
type Params struct {
	fx.In
	Propagator     otel.TextMapPropagator
	SpanProcessor  sdktrace.SpanProcessor
	Exporter       export.SpanExporter
	TracerProvider trace.TracerProvider
	Pusher         *push.Controller
	MetricPrivder  metric.MeterProvider
}

//StartTracer starts the distributed tracing service
func StartTracer(p Params, lifecycle fx.Lifecycle) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			global.SetTextMapPropagator(p.Propagator)
			global.SetTracerProvider(p.TracerProvider)
			global.SetMeterProvider(p.MetricPrivder)
			p.Pusher.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.SpanProcessor.Shutdown()
			err := p.Exporter.Shutdown(ctx)
			if err != nil {
				return err
			}
			p.Pusher.Stop() // pushes any last exports to the receiver
			return nil
		},
	})
}

//TracingModule provides a tracer that writes to a tracing backend
var TracingModule = fx.Options(fx.Provide(
	SamplingProvider,
	TraceConfigProvider,
	SpanProcessorProvider,
	TracerProvider,
	PushProvider,
	TextMapPropagatorProvider,
	ExporterProvider,
	TraceExporter,
	MetricExporter,
	MetricProvider))

//StdoutTracingModule provide a tracer that writes to stdout
var StdoutTracingModule = fx.Options(fx.Provide(
	SamplingProvider,
	TraceConfigProvider,
	SpanProcessorProvider,
	TracerProvider,
	PushProvider,
	TextMapPropagatorProvider,
	StdoutExporterProvider,
	StdoutTraceExporter,
	StdoutMetricExporter,
	MetricProvider))

/*
func main() {
	fx.New(
		tracing.TracingModule,
		fx.Invoke(tracing.StartTracer),
	).Run()

	tracer := global.Tracer("test-tracer")
	meter := global.Meter("test-meter")

	// labels represent additional key-value descriptors that can be bound to a
	// metric observer or recorder.
	commonLabels := []label.KeyValue{
		label.String("labelA", "chocolate"),
		label.String("labelB", "raspberry"),
		label.String("labelC", "vanilla"),
	}

	// Recorder metric example
	valuerecorder := metric.Must(meter).
		NewFloat64Counter(
			"an_important_metric",
			metric.WithDescription("Measures the cumulative epicness of the app"),
		).Bind(commonLabels...)
	defer valuerecorder.Unbind()

	// work begins
	ctx, span := tracer.Start(
		context.Background(),
		"CollectorExporter-Example",
		trace.WithAttributes(commonLabels...))
	defer span.End()
	for i := 0; i < 10; i++ {
		_, iSpan := tracer.Start(ctx, fmt.Sprintf("Sample-%d", i))
		log.Printf("Doing really hard work (%d / 10)\n", i+1)
		valuerecorder.Add(ctx, 1.0)

		<-time.After(time.Second)
		iSpan.End()
	}

	log.Printf("Done!")
}
*/

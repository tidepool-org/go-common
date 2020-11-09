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
	"log"
	"time"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/otlp"

	"go.opentelemetry.io/otel/propagators"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

//TraceConfig defines constants for configuring the trace system
type TraceConfig struct {
	Collector    string `envconfig:"OTEL_COLLECTOR_HOST" required:"true"`
	PodName      string `envconfig:"POD_NAME" required:"true"`
	PodNamespace string `envconfig:"POD_NAMESPACE" required:"true"`
	PodIP        string `envconfig:"POD_IP" required:"true"`
}

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider(traceConfig TraceConfig, sampler sdktrace.Sampler) func() {

	exp, err := otlp.NewExporter(
		otlp.WithInsecure(),
		otlp.WithAddress(traceConfig.Collector),
		otlp.WithGRPCDialOption(grpc.WithBlock()), // useful for testing
	)
	handleErr(err, "failed to create exporter")

	res := resource.New(
		semconv.ServiceNameKey.String(traceConfig.PodName),
		semconv.K8SNamespaceNameKey.String(traceConfig.PodNamespace),
		semconv.ServiceInstanceIDKey.String(traceConfig.PodIP),
	)

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	pusher := push.New(
		basic.New(
			simple.NewWithExactDistribution(),
			exp,
		),
		exp,
		push.WithPeriod(2*time.Second),
	)

	// set global propagator to tracecontext and baggage
	global.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{}))
	global.SetTracerProvider(tracerProvider)
	global.SetMeterProvider(pusher.MeterProvider())
	pusher.Start()

	return func() {
		bsp.Shutdown() // shutdown the processor
		handleErr(exp.Shutdown(context.Background()), "failed to stop exporter")
		pusher.Stop() // pushes any last exports to the receiver
	}
}

func handleErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

/*
func main() {
	log.Printf("Waiting for connection...")

	var traceConfig TraceConfig

	if err := envconfig.Process("", &traceConfig); err != nil {
		log.Fatalf("could not load trace config %v", err)
	}

	shutdown := initProvider(traceConfig, sdktrace.AlwaysSample())
	defer shutdown()

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

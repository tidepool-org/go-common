module github.com/tidepool-org/go-common

go 1.15

require (
	github.com/Shopify/sarama v1.27.0
	github.com/cloudevents/sdk-go/protocol/kafka_sarama/v2 v2.2.0
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c
	go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama v0.13.0
	go.opentelemetry.io/otel v0.13.0
	go.opentelemetry.io/otel/exporters/otlp v0.13.0
	go.opentelemetry.io/otel/sdk v0.13.0
	google.golang.org/grpc v1.32.0
)

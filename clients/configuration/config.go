package configuration

import (
	"crypto/tls"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// OutboundConfig contains how to communicate with the dependent services
type OutboundConfig struct {
	Protocol                string `default:"https"`
	ServerSecret            string `split_words:"true" required:"true"`
	AuthClientAddress       string `split_words:"true" required:"true"`
	PermissionClientAddress string `split_words:"true" required:"true"`
	MetricsClientAddress    string `split_words:"true" required:"true"`
	SeagullClientAddress    string `split_words:"true" required:"true"`
}

//OutboundConfigProvider provide outbound addresses
func OutboundConfigProvider() (OutboundConfig, error) {
	var config OutboundConfig
	err := envconfig.Process("tidepool", &config)
	if err != nil {
		return OutboundConfig{}, err
	}
	return config, nil
}

//InboundConfig describes how to receive inbound communication
type InboundConfig struct {
	Protocol      string `default:"http"`
	SslKeyFile    string `split_words:"true" default:""`
	SslCertFile   string `split_words:"true" default:""`
	ListenAddress string `split_words:"true" required:"true"`
}

// InboundConfigProvider provide inbound addresses
func InboundConfigProvider() (InboundConfig, error) {
	var config InboundConfig
	err := envconfig.Process("tidepool", &config)
	if err != nil {
		return InboundConfig{}, err
	}
	return config, nil
}

func httpClientProvider() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &http.Client{Transport: otelhttp.NewTransport(tr)}
}

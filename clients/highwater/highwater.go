package highwater

import (
	"github.com/tidepool-org/go-common/clients/disc"
	"net/http"
	"net/url"
)

//Generic client interface that we will implement and mock
type Client interface {
	postServer(eventName, token string, params map[string]string)
	postThisUser(eventName, token string, params map[string]string)
	postWithUser(userId, eventName, token string, params map[string]string)
	getHost() *url.URL
}

type HighwaterClient struct {
	hostGetter disc.HostGetter
	config     *HighwaterClientConfig
	httpClient *http.Client
}

type HighwaterClientBuilder struct {
	hostGetter disc.HostGetter
	config     *HighwaterClientConfig
	httpClient *http.Client
}

type HighwaterClientConfig struct {
	Name   string `json:"name"`   // The name of this server for use in obtaining a server token
	Secret string `json:"secret"` // The secret used along with the name to obtain a server token
}

func NewHighwaterClientBuilder() *HighwaterClientBuilder {
	return &HighwaterClientBuilder{
		config: &HighwaterClientConfig{},
	}
}

func (b *HighwaterClientBuilder) WithHostGetter(val disc.HostGetter) *HighwaterClientBuilder {
	b.hostGetter = val
	return b
}

func (b *HighwaterClientBuilder) WithHttpClient(val *http.Client) *HighwaterClientBuilder {
	b.httpClient = val
	return b
}

func (b *HighwaterClientBuilder) WithName(val string) *HighwaterClientBuilder {
	b.config.Name = val
	return b
}

func (b *HighwaterClientBuilder) WithSecret(val string) *HighwaterClientBuilder {
	b.config.Secret = val
	return b
}

func (b *HighwaterClientBuilder) WithConfig(val *HighwaterClientConfig) *HighwaterClientBuilder {
	return b.WithName(val.Name).WithSecret(val.Secret)
}

func (b *HighwaterClientBuilder) Build() *HighwaterClient {
	if b.hostGetter == nil {
		panic("HighwaterClient requires a hostGetter to be set")
	}
	if b.config.Name == "" {
		panic("HighwaterClient requires a name to be set")
	}
	if b.config.Secret == "" {
		panic("HighwaterClient requires a secret to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &HighwaterClient{
		hostGetter: b.hostGetter,
		httpClient: b.httpClient,
		config:     b.config,
	}
}

func (client *HighwaterClient) getHost() *url.URL {
	if hostArr := client.hostGetter.HostGet(); len(hostArr) > 0 {
		cpy := new(url.URL)
		*cpy = hostArr[0]
		return cpy
	} else {
		return nil
	}
}

func (client *HighwaterClient) postServer(eventName, token string, params map[string]string) {
	return
}

func (client *HighwaterClient) postThisUser(eventName, token string, params map[string]string) {
	return
}

func (client *HighwaterClient) postWithUser(userId, eventName, token string, params map[string]string) {
	return
}

package highwater

import (
	"bytes"
	"encoding/gob"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Client interface that we will implement and mock
type Client interface {
	PostServer(eventName, token string, params map[string]string)
	PostThisUser(eventName, token string, params map[string]string)
	PostWithUser(userId, eventName, token string, params map[string]string)
}

type HighwaterClient struct {
	host       url.URL
	config     *HighwaterClientConfig
	httpClient *http.Client
}

type HighwaterClientBuilder struct {
	host       *url.URL
	config     *HighwaterClientConfig
	httpClient *http.Client
	err        error
}

type HighwaterClientConfig struct {
	Name           string `json:"name"` // The name of this server for use in obtaining a server token
	MetricsSource  string `json:"metricsSource"`
	MetricsVersion string `json:"metricsVersion"`
}

func NewHighwaterClientBuilder() *HighwaterClientBuilder {
	return &HighwaterClientBuilder{
		config: &HighwaterClientConfig{},
	}
}

func (b *HighwaterClientBuilder) WithHost(host string) *HighwaterClientBuilder {
	h, err := url.Parse(host)
	if err != nil {
		b.err = err
	}
	b.host = h
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

func (b *HighwaterClientBuilder) WithSource(val string) *HighwaterClientBuilder {
	b.config.MetricsSource = val
	return b
}

func (b *HighwaterClientBuilder) WithVersion(val string) *HighwaterClientBuilder {
	b.config.MetricsVersion = val
	return b
}

func (b *HighwaterClientBuilder) WithConfig(val *HighwaterClientConfig) *HighwaterClientBuilder {
	return b.WithName(val.Name).WithSource(val.MetricsSource).WithVersion(val.MetricsVersion)
}

func (b *HighwaterClientBuilder) Build() *HighwaterClient {
	if b.err != nil {
		panic(b.err)
	}

	if b.host == nil {
		panic("HighwaterClient requires a host to be set")
	}
	if b.config.Name == "" {
		panic("HighwaterClient requires a name to be set")
	}
	if b.config.MetricsSource == "" {
		panic("HighwaterClient requires a source to be set")
	}

	if b.config.MetricsVersion == "" {
		panic("HighwaterClient requires a version to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &HighwaterClient{
		host:       *b.host,
		httpClient: b.httpClient,
		config:     b.config,
	}
}

func (client *HighwaterClient) getHost() url.URL {
	return client.host
}

func (client *HighwaterClient) adjustEventName(name string) string {
	src := client.config.MetricsSource
	src = strings.Replace(src, "-", " ", -1)

	return src + " - " + name
}

func (client *HighwaterClient) adjustEventParams(params map[string]string) []byte {
	params["sourceVersion"] = client.config.MetricsVersion

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(params); err != nil {
		log.Println("Error adjustEventParams ", err)
		return nil
	}
	return buf.Bytes()
}

func (client *HighwaterClient) PostServer(eventName, token string, params map[string]string) {
	host := client.getHost()
	host.Path = path.Join(host.Path, "server", client.config.Name, client.adjustEventName(eventName))

	req, err := http.NewRequest("GET", host.String(), bytes.NewBuffer(client.adjustEventParams(params)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-tidepool-session-token", token)

	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Error PostServer: [%s]  err[%v] ", req.URL, err)
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}

	return
}

func (client *HighwaterClient) PostThisUser(eventName, token string, params map[string]string) {
	host := client.getHost()
	host.Path = path.Join(host.Path, "thisuser", client.adjustEventName(eventName))

	req, err := http.NewRequest("GET", host.String(), bytes.NewBuffer(client.adjustEventParams(params)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-tidepool-session-token", token)

	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Error PostThisUser: [%s]  err[%v] ", req.URL, err)
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}

	return
}

func (client *HighwaterClient) PostWithUser(userId, eventName, token string, params map[string]string) {
	host := client.getHost()
	host.Path = path.Join(host.Path, "user", userId, client.adjustEventName(eventName))

	req, err := http.NewRequest("GET", host.String(), bytes.NewBuffer(client.adjustEventParams(params)))
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-tidepool-session-token", token)

	if _, err := client.httpClient.Do(req); err != nil {
		log.Printf("Error PostWithUser: [%s]  err[%v] ", req.URL, err)
	}

	return
}

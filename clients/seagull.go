package clients

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/tidepool-org/go-common/clients/configuration"
	"github.com/tidepool-org/go-common/clients/disc"
	"github.com/tidepool-org/go-common/clients/status"
	"go.uber.org/fx"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/semconv"
)

type (
	Seagull interface {
		// Retrieves arbitrary collection information from metadata
		//
		// userID -- the Tidepool-assigned userId
		// hashName -- the name of what we are trying to get
		// token -- a server token or the user token
		GetPrivatePair(ctx context.Context, userID, hashName, token string) *PrivatePair
		// Retrieves arbitrary collection information from metadata
		//
		// userID -- the Tidepool-assigned userId
		// collectionName -- the collection being retrieved
		// token -- a server token or the user token
		// v - the interface to return the value in
		GetCollection(ctx context.Context, userID, collectionName, token string, v interface{}) error
	}

	SeagullClient struct {
		httpClient *http.Client // store a reference to the http client so we can reuse it
		host       *url.URL
	}

	seagullClientBuilder struct {
		httpClient *http.Client
		host       *url.URL
	}

	PrivatePair struct {
		ID    string
		Value string
	}
)

//SeagullModule provides a Seagull client
var SeagullModule fx.Option = fx.Options(fx.Provide(SeagullProvider))

func SeagullProvider(config configuration.OutboundConfig, httpClient *http.Client) Seagull {
	host, _ := url.Parse(config.SeagullClientAddress)
	return NewSeagullClientBuilder().
		WithHost(host).
		WithHttpClient(httpClient).
		Build()
}

func NewSeagullClientBuilder() *seagullClientBuilder {
	return &seagullClientBuilder{}
}

func (b *seagullClientBuilder) WithHttpClient(httpClient *http.Client) *seagullClientBuilder {
	b.httpClient = httpClient
	return b
}

func (b *seagullClientBuilder) WithHostGetter(hostGetter disc.HostGetter) *seagullClientBuilder {
	b.host = &hostGetter.HostGet()[0]
	return b
}

func (b *seagullClientBuilder) WithHost(host *url.URL) *seagullClientBuilder {
	b.host = host
	return b
}

func (b *seagullClientBuilder) Build() *SeagullClient {
	if b.httpClient == nil {
		panic("seagullClient requires an httpClient to be set")
	}
	if b.host == nil {
		panic("seagullClient requires a hostGetter to be set")
	}
	return &SeagullClient{
		httpClient: b.httpClient,
		host:       b.host,
	}
}

func (client *SeagullClient) GetPrivatePair(ctx context.Context, userID, hashName, token string) *PrivatePair {
	host := client.host
	if host == nil {
		return nil
	}
	host.Path = path.Join(host.Path, userID, "private", hashName)

	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "GetPrivatePair", trace.WithAttributes(semconv.PeerServiceKey.String("seagull")))
	defer span.End()
	req, _ := http.NewRequestWithContext(spanCtx, "GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Problem when looking up private pair for userID[%s]. %s", userID, err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Unknown response code[%v] from service[%v]", res.StatusCode, req.URL)
		return nil
	}

	var retVal PrivatePair
	if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
		log.Println("Error parsing JSON results", err)
		return nil
	}
	return &retVal
}

func (client *SeagullClient) GetCollection(ctx context.Context, userID, collectionName, token string, v interface{}) error {
	host := client.host
	if host == nil {
		return nil
	}
	host.Path = path.Join(host.Path, userID, collectionName)
	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "GetCollection", trace.WithAttributes(semconv.PeerServiceKey.String("seagull")))
	defer span.End()

	req, _ := http.NewRequestWithContext(spanCtx, "GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Problem when looking up collection for userID[%s]. %s", userID, err)
		return err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
			log.Println("Error parsing JSON results", err)
			return err
		}
		return nil
	case http.StatusNotFound:
		log.Printf("No [%s] collection found for [%s]", collectionName, userID)
		return nil
	default:
		return &status.StatusError{Status: status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}

}

package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/tidepool-org/go-common/clients/configuration"
	"github.com/tidepool-org/go-common/clients/disc"
	"github.com/tidepool-org/go-common/clients/shoreline"
	"github.com/tidepool-org/go-common/clients/status"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/semconv"
	"go.uber.org/fx"
)

type (
	//Gatekeeper  Interface so that we can mock gatekeeperClient for tests
	Gatekeeper interface {
		//userID  -- the Tidepool-assigned userID
		//groupID  -- the Tidepool-assigned groupID
		//
		// returns the Permissions
		UserInGroup(ctx context.Context, userID, groupID string) (Permissions, error)

		//groupID  -- the Tidepool-assigned groupID
		//
		// returns the map of user id to Permissions
		UsersInGroup(ctx context.Context, groupID string) (UsersPermissions, error)

		//userID  -- the Tidepool-assigned userID
		//groupID  -- the Tidepool-assigned groupID
		//permissions -- the permisson we want to give the user for the group
		SetPermissions(ctx context.Context, userID, groupID string, permissions Permissions) (Permissions, error)
	}

	GatekeeperClient struct {
		httpClient    *http.Client  // store a reference to the http client so we can reuse it
		tokenProvider TokenProvider // An object that provides tokens for communicating with gatekeeper
		host          *url.URL
	}

	gatekeeperClientBuilder struct {
		httpClient    *http.Client  // store a reference to the http client so we can reuse it
		tokenProvider TokenProvider // An object that provides tokens for communicating with gatekeeper
		host          *url.URL
	}

	Permission       map[string]interface{}
	Permissions      map[string]Permission
	UsersPermissions map[string]Permissions
)

var Allowed Permission = Permission{}

//GatekeeperModule provides a Gatekeeper client
var GatekeeperModule fx.Option = fx.Options(fx.Provide(Provider))

//Provider creats a Gatekeeper client
func Provider(config configuration.OutboundConfig, shoreline shoreline.Client, httpClient *http.Client) Gatekeeper {
	host, _ := url.Parse(config.PermissionClientAddress)
	return NewGatekeeperClientBuilder().
		WithHost(host).
		WithHttpClient(httpClient).
		WithTokenProvider(shoreline).
		Build()
}

func NewGatekeeperClientBuilder() *gatekeeperClientBuilder {
	return &gatekeeperClientBuilder{}
}

func (b *gatekeeperClientBuilder) WithHttpClient(httpClient *http.Client) *gatekeeperClientBuilder {
	b.httpClient = httpClient
	return b
}

func (b *gatekeeperClientBuilder) WithHostGetter(hostGetter disc.HostGetter) *gatekeeperClientBuilder {
	b.host = &hostGetter.HostGet()[0]
	return b
}

func (b *gatekeeperClientBuilder) WithHost(host *url.URL) *gatekeeperClientBuilder {
	b.host = host
	return b
}

func (b *gatekeeperClientBuilder) WithTokenProvider(tokenProvider TokenProvider) *gatekeeperClientBuilder {
	b.tokenProvider = tokenProvider
	return b
}

func (b *gatekeeperClientBuilder) Build() *GatekeeperClient {
	if b.host == nil {
		panic("gatekeeperClient requires a host to be set")
	}
	if b.tokenProvider == nil {
		panic("gatekeeperClient requires a tokenProvider to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &GatekeeperClient{
		httpClient:    b.httpClient,
		host:          b.host,
		tokenProvider: b.tokenProvider,
	}
}

func (client *GatekeeperClient) UserInGroup(ctx context.Context, userID, groupID string) (Permissions, error) {
	host := client.host
	host.Path = path.Join(host.Path, "access", groupID, userID)

	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "UserInGroup", trace.WithAttributes(semconv.PeerServiceKey.String("gatekeeper")))
	defer span.End()

	req, _ := http.NewRequestWithContext(spanCtx, "GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide(ctx))

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		retVal := make(Permissions)
		if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
			log.Println(err)
			return nil, &status.StatusError{Status: status.NewStatus(500, "UserInGroup Unable to parse response.")}
		}
		return retVal, nil
	} else if res.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, &status.StatusError{Status: status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

func (client *GatekeeperClient) UsersInGroup(ctx context.Context, groupID string) (UsersPermissions, error) {
	host := client.host
	host.Path = path.Join(host.Path, "access", groupID)

	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "UserInGroup", trace.WithAttributes(semconv.PeerServiceKey.String("gatekeeper")))
	defer span.End()

	req, _ := http.NewRequestWithContext(spanCtx, "GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide(ctx))

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		retVal := make(UsersPermissions)
		if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
			log.Println(err)
			return nil, &status.StatusError{Status: status.NewStatus(500, "UserInGroup Unable to parse response.")}
		}
		return retVal, nil
	} else if res.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, &status.StatusError{Status: status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

func (client *GatekeeperClient) SetPermissions(ctx context.Context, userID, groupID string, permissions Permissions) (Permissions, error) {
	host := client.host
	host.Path = path.Join(host.Path, "access", groupID, userID)

	if jsonPerms, err := json.Marshal(permissions); err != nil {
		log.Println(err)
		return nil, &status.StatusError{Status: status.NewStatusf(http.StatusInternalServerError, "Error marshaling the permissons [%s]", err)}
	} else {
		tr := global.Tracer("go-common tracer")

		spanCtx, span := tr.Start(ctx, "UserInGroup", trace.WithAttributes(semconv.PeerServiceKey.String("gatekeeper")))
		defer span.End()

		req, _ := http.NewRequestWithContext(spanCtx, "POST", host.String(), bytes.NewBuffer(jsonPerms))
		req.Header.Set("content-type", "application/json")
		req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide(ctx))

		res, err := client.httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == 200 {
			retVal := make(Permissions)
			if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
				log.Printf("SetPermissions: Unable to parse response: [%s]", err.Error())
				return nil, &status.StatusError{Status: status.NewStatus(500, "SetPermissions: Unable to parse response:")}
			}
			return retVal, nil
		} else if res.StatusCode == 404 {
			return nil, nil
		} else {
			return nil, &status.StatusError{Status: status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
		}
	}
}

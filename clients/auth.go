package clients

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/tidepool-org/go-common/clients/disc"
	"github.com/tidepool-org/go-common/clients/status"
	"github.com/tidepool-org/go-common/errors"
)

type (
	//Interface so that we can mock authClient for tests
	Auth interface {
		//userID  -- the Tidepool-assigned userID
		//
		// returns the Auth Sources for the user
		CreateRestrictedToken(userID string, expirationTime time.Time, paths []string, token string) (*RestrictedToken, error)
	}

	AuthClient struct {
		httpClient    *http.Client    // store a reference to the http client so we can reuse it
		hostGetter    disc.HostGetter // The getter that provides the host to talk to for the client
		tokenProvider TokenProvider   // An object that provides tokens for communicating with data
	}

	authClientBuilder struct {
		httpClient    *http.Client    // store a reference to the http client so we can reuse it
		hostGetter    disc.HostGetter // The getter that provides the host to talk to for the client
		tokenProvider TokenProvider   // An object that provides tokens for communicating with data
	}
)

// RestrictedToken is the data structure returned from a successful create restricted token query.
type RestrictedToken struct {
	ID             string     `json:"id" bson:"id"`
	UserID         string     `json:"userId" bson:"userId"`
	Paths          *[]string  `json:"paths,omitempty" bson:"paths,omitempty"`
	ExpirationTime time.Time  `json:"expirationTime" bson:"expirationTime"`
	CreatedTime    time.Time  `json:"createdTime" bson:"createdTime"`
	ModifiedTime   *time.Time `json:"modifiedTime,omitempty" bson:"modifiedTime,omitempty"`
}

func NewAuthClientBuilder() *authClientBuilder {
	return &authClientBuilder{}
}

func (b *authClientBuilder) WithHttpClient(httpClient *http.Client) *authClientBuilder {
	b.httpClient = httpClient
	return b
}

func (b *authClientBuilder) WithHostGetter(hostGetter disc.HostGetter) *authClientBuilder {
	b.hostGetter = hostGetter
	return b
}

func (b *authClientBuilder) WithTokenProvider(tokenProvider TokenProvider) *authClientBuilder {
	b.tokenProvider = tokenProvider
	return b
}

func (b *authClientBuilder) Build() *AuthClient {
	if b.hostGetter == nil {
		panic("authClient requires a hostGetter to be set")
	}
	if b.tokenProvider == nil {
		panic("authClient requires a tokenProvider to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &AuthClient{
		httpClient:    b.httpClient,
		hostGetter:    b.hostGetter,
		tokenProvider: b.tokenProvider,
	}
}

// CreateRestrictedToken creates a restricted token for a given user
func (client *AuthClient) CreateRestrictedToken(userID string, expirationTime time.Time, paths []string, token string) (*RestrictedToken, error) {
	host := client.getHost()
	if host == nil {
		return nil, errors.New("No known user-api hosts.")
	}

	host.Path = path.Join(host.Path, "v1", "users", userID, "restricted_tokens")

	req, _ := http.NewRequest("POST", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create restricted token")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusCreated:
		var td RestrictedToken
		if err = json.NewDecoder(res.Body).Decode(&td); err != nil {
			log.Println("Error parsing JSON results", err)
			return nil, nil
		}
		return &td, nil
	default:
		return nil, &status.StatusError{
			Status: status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL),
		}
	}
}

func (client *AuthClient) getHost() *url.URL {
	if hostArr := client.hostGetter.HostGet(); len(hostArr) > 0 {
		cpy := new(url.URL)
		*cpy = hostArr[0]
		return cpy
	} else {
		return nil
	}
}

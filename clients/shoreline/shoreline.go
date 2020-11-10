// Package shoreline is a client module to support server-side use of the Tidepool
// service called user-api.
package shoreline

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/tidepool-org/go-common/clients/disc"
	"github.com/tidepool-org/go-common/clients/status"
	"github.com/tidepool-org/go-common/errors"
	"github.com/tidepool-org/go-common/jepson"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/semconv"
)

// Client interface that we will implement and mock
type Client interface {
	Start(ctx context.Context) error
	Close(ctx context.Context)
	Login(ctx context.Context, username, password string) (*UserData, string, error)
	Signup(ctx context.Context, username, password, email string) (*UserData, error)
	CheckToken(ctx context.Context, token string) *TokenData
	TokenProvide(ctx context.Context) string
	GetUser(ctx context.Context, userID, token string) (*UserData, error)
	UpdateUser(ctx context.Context, userID string, userUpdate UserUpdate, token string) error
}

// ShorelineClient manages the local data for a client. A client is intended to be shared among multiple
// goroutines so it's OK to treat it as a singleton (and probably a good idea).
type ShorelineClient struct {
	httpClient *http.Client // store a reference to the http client so we can reuse it
	host       *url.URL
	config     *ShorelineClientConfig // Configuration for the client

	mut         sync.Mutex
	serverToken string         // stores the most recently received server token
	closed      chan chan bool // Channel to communicate that the object has been closed
}

type ShorelineClientConfig struct {
	Name                 string          `json:"name"`                 // The name of this server for use in obtaining a server token
	Secret               string          `json:"secret"`               // The secret used along with the name to obtain a server token
	TokenRefreshInterval jepson.Duration `json:"tokenRefreshInterval"` // The amount of time between refreshes of the server token
}

// UserData is the data structure returned from a successful Login query.
type UserData struct {
	UserID         string   `json:"userid,omitempty"`         // the tidepool-assigned user ID
	Username       string   `json:"username,omitempty"`       // the user-assigned name for the login (usually an email address)
	Emails         []string `json:"emails,omitempty"`         // the array of email addresses associated with this account
	PasswordExists bool     `json:"passwordExists,omitempty"` // Does a password exist for the user?
	Roles          []string `json:"roles,omitempty"`          // User roles
	EmailVerified  bool     `json:"emailVerified,omitempty"`  // the user has verified the email used as part of signup
	TermsAccepted  string   `json:"termsAccepted,omitempty"`  // When were the terms accepted
}

// UserUpdate is the data structure for updating of a users details
type UserUpdate struct {
	Username      *string   `json:"username,omitempty"`
	Emails        *[]string `json:"emails,omitempty"`
	Password      *string   `json:"password,omitempty"`
	Roles         *[]string `json:"roles,omitempty"`
	EmailVerified *bool     `json:"emailVerified,omitempty"`
}

// TokenData is the data structure returned from a successful CheckToken query.
type TokenData struct {
	UserID   string // the UserID stored in the token
	IsServer bool   // true or false depending on whether the token was a servertoken
}

type ShorelineClientBuilder struct {
	host       *url.URL
	config     *ShorelineClientConfig
	httpClient *http.Client
}

func (u *UserData) IsCustodial() bool {
	return !u.PasswordExists
}

func (u *UserData) HasRole(role string) bool {
	for _, userRole := range u.Roles {
		if userRole == role {
			return true
		}
	}
	return false
}

func (u *UserData) IsClinic() bool {
	return u.HasRole("clinic")
}

func (u *UserUpdate) HasUpdates() bool {
	return u.Username != nil || u.Emails != nil || u.Password != nil || u.Roles != nil || u.EmailVerified != nil
}

func NewShorelineClientBuilder() *ShorelineClientBuilder {
	return &ShorelineClientBuilder{
		config: &ShorelineClientConfig{
			TokenRefreshInterval: jepson.Duration(6 * time.Hour),
		},
	}
}

func (b *ShorelineClientBuilder) WithHost(host *url.URL) *ShorelineClientBuilder {
	b.host = host
	return b
}

func (b *ShorelineClientBuilder) WithHostGetter(val disc.HostGetter) *ShorelineClientBuilder {
	b.host = &val.HostGet()[0]
	return b
}

func (b *ShorelineClientBuilder) WithHttpClient(val *http.Client) *ShorelineClientBuilder {
	b.httpClient = val
	return b
}

func (b *ShorelineClientBuilder) WithName(val string) *ShorelineClientBuilder {
	b.config.Name = val
	return b
}

func (b *ShorelineClientBuilder) WithSecret(val string) *ShorelineClientBuilder {
	b.config.Secret = val
	return b
}

func (b *ShorelineClientBuilder) WithTokenRefreshInterval(val time.Duration) *ShorelineClientBuilder {
	b.config.TokenRefreshInterval = jepson.Duration(val)
	return b
}

func (b *ShorelineClientBuilder) WithConfig(val *ShorelineClientConfig) *ShorelineClientBuilder {
	return b.WithName(val.Name).WithSecret(val.Secret).WithTokenRefreshInterval(time.Duration(val.TokenRefreshInterval))
}

func (b *ShorelineClientBuilder) Build() *ShorelineClient {
	if b.host == nil {
		panic("shorelineClient requires a host to be set")
	}
	if b.config.Name == "" {
		panic("shorelineClient requires a name to be set")
	}
	if b.config.Secret == "" {
		panic("shorelineClient requires a secret to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &ShorelineClient{
		host:       b.host,
		httpClient: b.httpClient,
		config:     b.config,

		closed: make(chan chan bool),
	}
}

// Start starts the client and makes it ready for us.  This must be done before using any of the functionality
// that requires a server token
func (client *ShorelineClient) Start(ctx context.Context) error {
	if err := client.serverLogin(ctx); err != nil {
		log.Printf("Problem with initial server token acquisition, [%v]", err)
		panic(err)
	}

	go func() {
		for {
			timer := time.After(time.Duration(client.config.TokenRefreshInterval))
			select {
			case twoWay := <-client.closed:
				twoWay <- true
				return
			case <-timer:
				if err := client.serverLogin(ctx); err != nil {
					log.Print("Error when refreshing server login", err)
					panic(err)
				}
			}
		}
	}()
	return nil
}

func (client *ShorelineClient) Close(ctx context.Context) {
	twoWay := make(chan bool)
	client.closed <- twoWay
	<-twoWay

	client.mut.Lock()
	defer client.mut.Unlock()
	client.serverToken = ""
}

// serverLogin issues a request to the server for a login, using the stored
// secret that was passed in on the creation of the client object. If
// successful, it stores the returned token in ServerToken.
func (client *ShorelineClient) serverLogin(ctx context.Context) error {
	host := client.host
	if host == nil {
		return errors.New("No known user-api hosts")
	}

	host.Path = path.Join(host.Path, "serverlogin")

	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "serverLogin", trace.WithAttributes(semconv.PeerServiceKey.String("shoreline")))
	defer span.End()

	req, _ := http.NewRequestWithContext(spanCtx, "POST", host.String(), nil)
	req.Header.Add("x-tidepool-server-name", client.config.Name)
	req.Header.Add("x-tidepool-server-secret", client.config.Secret)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failure to obtain a server token")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return &status.StatusError{
			status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
	token := res.Header.Get("x-tidepool-session-token")

	client.mut.Lock()
	defer client.mut.Unlock()
	client.serverToken = token

	return nil
}

func extractUserData(r io.Reader) (*UserData, error) {
	var ud UserData
	if err := json.NewDecoder(r).Decode(&ud); err != nil {
		return nil, err
	}
	return &ud, nil
}

// Signs up a new platfrom user
// Returns a UserData object if successful
func (client *ShorelineClient) Signup(ctx context.Context, username, password, email string) (*UserData, error) {
	host := client.host
	if host == nil {
		return nil, errors.New("No known user-api hosts.")
	}

	host.Path = path.Join(host.Path, "user")
	data := []byte(fmt.Sprintf(`{"username": "%s", "password": "%s","emails":["%s"]}`, username, password, email))

	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "Signup", trace.WithAttributes(semconv.PeerServiceKey.String("shoreline")))
	defer span.End()

	req, _ := http.NewRequestWithContext(spanCtx, "POST", host.String(), bytes.NewBuffer(data))
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusCreated:
		ud, err := extractUserData(res.Body)
		if err != nil {
			return nil, err
		}

		return ud, nil
	default:
		return nil, &status.StatusError{status.NewStatus(res.StatusCode, "There was an issue trying to signup a new user")}
	}
}

// Login logs in a user with a username and password. Returns a UserData object if successful
// and also stores the returned login token into ClientToken.
func (client *ShorelineClient) Login(ctx context.Context, username, password string) (*UserData, string, error) {
	host := client.host
	if host == nil {
		return nil, "", errors.New("No known user-api hosts.")
	}

	host.Path = path.Join(host.Path, "login")

	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "Login", trace.WithAttributes(semconv.PeerServiceKey.String("shoreline")))
	defer span.End()

	req, _ := http.NewRequestWithContext(spanCtx, "POST", host.String(), nil)
	req.SetBasicAuth(username, password)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		ud, err := extractUserData(res.Body)
		if err != nil {
			return nil, "", err
		}

		return ud, res.Header.Get("x-tidepool-session-token"), nil
	case 404:
		return nil, "", nil
	default:
		return nil, "", &status.StatusError{
			status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

// CheckToken tests a token with the user-api to make sure it's current;
// if so, it returns the data encoded in the token.
func (client *ShorelineClient) CheckToken(ctx context.Context, token string) *TokenData {
	host := client.host
	if host == nil {
		return nil
	}

	host.Path = path.Join(host.Path, "token", token)

	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "CheckToken", trace.WithAttributes(semconv.PeerServiceKey.String("shoreline")))
	defer span.End()

	req, _ := http.NewRequestWithContext(spanCtx, "GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.serverToken)

	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Println("Error checking token", err)
		return nil
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		var td TokenData
		if err = json.NewDecoder(res.Body).Decode(&td); err != nil {
			log.Println("Error parsing JSON results", err)
			return nil
		}
		return &td
	case 404:
		return nil
	default:
		log.Printf("Unknown response code[%d] from service[%s]", res.StatusCode, req.URL)
		return nil
	}
}

func (client *ShorelineClient) TokenProvide(ctx context.Context) string {
	client.mut.Lock()
	defer client.mut.Unlock()

	return client.serverToken
}

// Get user details for the given user
// In this case the userID could be the actual ID or an email address
func (client *ShorelineClient) GetUser(ctx context.Context, userID, token string) (*UserData, error) {
	host := client.host
	if host == nil {
		return nil, errors.New("No known user-api hosts.")
	}

	host.Path = path.Join(host.Path, "user", userID)

	tr := global.Tracer("go-common tracer")

	spanCtx, span := tr.Start(ctx, "GetUser", trace.WithAttributes(semconv.PeerServiceKey.String("shoreline")))
	defer span.End()

	req, _ := http.NewRequestWithContext(spanCtx, "GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failure to get a user")
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		ud, err := extractUserData(res.Body)
		if err != nil {
			return nil, err
		}
		return ud, nil
	case http.StatusNoContent:
		return &UserData{}, nil
	default:
		return nil, &status.StatusError{
			status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

// Get user details for the given user
// In this case the userID could be the actual ID or an email address
func (client *ShorelineClient) UpdateUser(ctx context.Context, userID string, userUpdate UserUpdate, token string) error {
	host := client.host
	if host == nil {
		return errors.New("No known user-api hosts.")
	}

	//structure that the update are given to us in
	type updatesToApply struct {
		Updates UserUpdate `json:"updates"`
	}

	host.Path = path.Join(host.Path, "user", userID)

	if jsonUser, err := json.Marshal(updatesToApply{Updates: userUpdate}); err != nil {
		return &status.StatusError{
			status.NewStatusf(http.StatusInternalServerError, "Error getting user updates [%s]", err.Error())}
	} else {

		tr := global.Tracer("go-common tracer")

		spanCtx, span := tr.Start(ctx, "UpdateUser", trace.WithAttributes(semconv.PeerServiceKey.String("shoreline")))
		defer span.End()

		req, _ := http.NewRequestWithContext(spanCtx, "PUT", host.String(), bytes.NewBuffer(jsonUser))
		req.Header.Add("x-tidepool-session-token", token)

		res, err := client.httpClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "Failure to get a user")
		}
		defer res.Body.Close()

		switch res.StatusCode {
		case http.StatusOK:
			return nil
		default:
			return &status.StatusError{
				status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
		}
	}
}

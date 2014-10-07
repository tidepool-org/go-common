// This is a client module to support server-side use of the Tidepool
// service called user-api.
package shoreline

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/tidepool-org/go-common/clients/disc"
	"github.com/tidepool-org/go-common/clients/status"
	"github.com/tidepool-org/go-common/errors"
	"github.com/tidepool-org/go-common/jepson"
)

const (
	X_TIDEPOOL_SERVER_NAME   = "x-tidepool-server-name"
	X_TIDEPOOL_SERVER_SECRET = "x-tidepool-server-secret"
	X_TIDEPOOL_SESSION_TOKEN = "x-tidepool-session-token"
)

type (

	//Generic client interface that we will implement and mock
	Client interface {
		Close()
		CheckToken(token string) *TokenData
		getHost() *url.URL
		GetUser(userID, token string) (*UserData, error)
		serverLogin() error
		Start() error
		TokenProvide() string
	}

	// UserApiClient manages the local data for a client. A client is intended to be shared among multiple
	// goroutines so it's OK to treat it as a singleton (and probably a good idea).
	ShorelineClient struct {
		httpClient *http.Client           // store a reference to the http client so we can reuse it
		hostGetter disc.HostGetter        // The getter that provides the host to talk to for the client
		config     *ShorelineClientConfig // Configuration for the client

		mut         sync.Mutex
		serverToken string         // stores the most recently received server token
		closed      chan chan bool // Channel to communicate that the object has been closed
	}

	ShorelineClientConfig struct {
		Name                 string          `json:"name"`                 // The name of this server for use in obtaining a server token
		Secret               string          `json:"secret"`               // The secret used along with the name to obtain a server token
		TokenRefreshInterval jepson.Duration `json:"tokenRefreshInterval"` // The amount of time between refreshes of the server token
	}

	// UserData is the data structure returned from a successful Login query.
	UserData struct {
		UserID   string   // the tidepool-assigned user ID
		UserName string   // the user-assigned name for the login (usually an email address)
		Emails   []string // the array of email addresses associated with this account
	}

	// TokenData is the data structure returned from a successful CheckToken query.
	TokenData struct {
		UserID   string // the UserID stored in the token
		IsServer bool   // true or false depending on whether the token was a servertoken
	}

	ShorelineClientBuilder struct {
		hostGetter disc.HostGetter
		config     *ShorelineClientConfig
		httpClient *http.Client
	}
)

func NewShorelineClientBuilder() *ShorelineClientBuilder {
	return &ShorelineClientBuilder{
		config: &ShorelineClientConfig{
			TokenRefreshInterval: jepson.Duration(6 * time.Hour),
		},
	}
}

func (b *ShorelineClientBuilder) WithHostGetter(val disc.HostGetter) *ShorelineClientBuilder {
	b.hostGetter = val
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
	if b.hostGetter == nil {
		panic("shorelineClient requires a hostGetter to be set")
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
		hostGetter: b.hostGetter,
		httpClient: b.httpClient,
		config:     b.config,

		closed: make(chan chan bool),
	}
}

// Start starts the client and makes it ready for us.  This must be done before using any of the functionality
// that requires a server token
func (client *ShorelineClient) Start() error {
	if err := client.serverLogin(); err != nil {
		log.Printf("Problem with initial server token acquisition, [%v]", err)
	}

	go func() {
		for {
			timer := time.After(time.Duration(client.config.TokenRefreshInterval))
			select {
			case twoWay := <-client.closed:
				twoWay <- true
				return
			case <-timer:
				if err := client.serverLogin(); err != nil {
					log.Print("Error when refreshing server login", err)
				}
			}
		}
	}()
	return nil
}

func (client *ShorelineClient) Close() {
	twoWay := make(chan bool)
	client.closed <- twoWay
	<-twoWay

	client.mut.Lock()
	defer client.mut.Unlock()
	client.serverToken = ""
}

// Get user details for the given user
// In this case the userID could be the actual ID or an email address
func (client *ShorelineClient) GetUser(userID, token string) (*UserData, error) {
	host := client.getHost()
	if host == nil {
		return nil, errors.New("No known user-api hosts.")
	}

	host.Path += fmt.Sprintf("user/%s", userID)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add(X_TIDEPOOL_SESSION_TOKEN, token)

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

// serverLogin issues a request to the server for a login, using the stored
// secret that was passed in on the creation of the client object. If
// successful, it stores the returned token in ServerToken.
func (client *ShorelineClient) serverLogin() error {
	host := client.getHost()
	if host == nil {
		return errors.New("No known user-api hosts.")
	}

	host.Path += "/serverlogin"

	req, _ := http.NewRequest("POST", host.String(), nil)
	req.Header.Add(X_TIDEPOOL_SERVER_NAME, client.config.Name)
	req.Header.Add(X_TIDEPOOL_SERVER_SECRET, client.config.Secret)

	res, err := client.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "Failure to obtain a server token")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return &status.StatusError{
			status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
	token := res.Header.Get(X_TIDEPOOL_SESSION_TOKEN)

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

// Login logs in a user with a username and password. Returns a UserData object if successful
// and also stores the returned login token into ClientToken.
func (client *ShorelineClient) Login(username, password string) (*UserData, string, error) {
	host := client.getHost()
	if host == nil {
		return nil, "", errors.New("No known user-api hosts.")
	}

	host.Path += "/login"

	req, _ := http.NewRequest("POST", host.String(), nil)
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

		return ud, res.Header.Get(X_TIDEPOOL_SESSION_TOKEN), nil
	case 404:
		return nil, "", nil
	default:
		return nil, "", &status.StatusError{
			status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

// CheckToken tests a token with the user-api to make sure it's current;
// if so, it returns the data encoded in the token.
func (client *ShorelineClient) CheckToken(token string) *TokenData {
	host := client.getHost()
	if host == nil {
		return nil
	}

	host.Path += "/token/" + token

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add(X_TIDEPOOL_SESSION_TOKEN, client.serverToken)

	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Println("Error checking token", err)
		return nil
	}

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

func (client *ShorelineClient) TokenProvide() string {
	client.mut.Lock()
	defer client.mut.Unlock()

	return client.serverToken
}

func (client *ShorelineClient) getHost() *url.URL {
	if hostArr := client.hostGetter.HostGet(); len(hostArr) > 0 {
		cpy := new(url.URL)
		*cpy = hostArr[0]
		return cpy
	} else {
		return nil
	}
}

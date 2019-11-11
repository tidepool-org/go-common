package clients

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/tidepool-org/go-common/clients/status"
)

type (
	//Gatekeeper Inteface so that we can mock gatekeeperClient for tests
	Gatekeeper interface {
		//userID  -- the Tidepool-assigned userID
		//groupID  -- the Tidepool-assigned groupID
		//
		// returns the Permissions
		UserInGroup(userID, groupID string) (Permissions, error)

		//groupID  -- the Tidepool-assigned groupID
		//
		// returns the map of user id to Permissions
		UsersInGroup(groupID string) (UsersPermissions, error)

		//userID  -- the Tidepool-assigned userID
		//groupID  -- the Tidepool-assigned groupID
		//permissions -- the permisson we want to give the user for the group
		SetPermissions(userID, groupID string, permissions Permissions) (Permissions, error)
	}

	gatekeeperClient struct {
		httpClient    *http.Client  // store a reference to the http client so we can reuse it
		host          url.URL       // The getter that provides the host to talk to for the client
		tokenProvider TokenProvider // An object that provides tokens for communicating with gatekeeper
	}

	gatekeeperClientBuilder struct {
		httpClient    *http.Client  // store a reference to the http client so we can reuse it
		host          *url.URL      // The getter that provides the host to talk to for the client
		tokenProvider TokenProvider // An object that provides tokens for communicating with gatekeeper
		err           error
	}

	Permission       map[string]interface{}
	Permissions      map[string]Permission
	UsersPermissions map[string]Permissions
)

var (
	Allowed Permission = Permission{}
)

func NewGatekeeperClientBuilder() *gatekeeperClientBuilder {
	return &gatekeeperClientBuilder{}
}

func (b *gatekeeperClientBuilder) WithHttpClient(httpClient *http.Client) *gatekeeperClientBuilder {
	b.httpClient = httpClient
	return b
}

func (b *gatekeeperClientBuilder) WithHost(host string) *gatekeeperClientBuilder {
	h, err := url.Parse(host)
	if err != nil {
		b.err = err
	}
	b.host = h
	return b
}

func (b *gatekeeperClientBuilder) WithTokenProvider(tokenProvider TokenProvider) *gatekeeperClientBuilder {
	b.tokenProvider = tokenProvider
	return b
}

func (b *gatekeeperClientBuilder) Build() *gatekeeperClient {
	if b.err != nil {
		panic(b.err)
	}
	if b.host == nil {
		panic("gatekeeperClient requires a host to be set")
	}
	if b.tokenProvider == nil {
		panic("gatekeeperClient requires a tokenProvider to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &gatekeeperClient{
		httpClient:    b.httpClient,
		host:          *b.host,
		tokenProvider: b.tokenProvider,
	}
}

func (client *gatekeeperClient) UserInGroup(userID, groupID string) (Permissions, error) {
	host := client.getHost()
	host.Path = path.Join(host.Path, "access", groupID, userID)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

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

func (client *gatekeeperClient) UsersInGroup(groupID string) (UsersPermissions, error) {
	host := client.getHost()
	host.Path = path.Join(host.Path, "access", groupID)

	req, err := http.NewRequest("GET", host.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

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

func (client *gatekeeperClient) SetPermissions(userID, groupID string, permissions Permissions) (Permissions, error) {
	host := client.getHost()
	host.Path = path.Join(host.Path, "access", groupID, userID)

	jsonPerms, err := json.Marshal(permissions)

	if err != nil {
		log.Println(err)
		return nil, &status.StatusError{Status: status.NewStatusf(http.StatusInternalServerError, "Error marshaling the permissions [%s]", err)}
	}

	req, err := http.NewRequest("POST", host.String(), bytes.NewBuffer(jsonPerms))
	if err != nil {
		panic(err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

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

func (client *gatekeeperClient) getHost() url.URL {
	return client.host
}

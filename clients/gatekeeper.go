package clients

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/tidepool-org/go-common/clients/disc"
	"github.com/tidepool-org/go-common/clients/status"
	"github.com/tidepool-org/go-common/errors"
)

type (
	//Inteface so that we can mock gatekeeperClient for tests
	Gatekeeper interface {
		UserInGroup(userID, groupID string) (map[string]Permissions, error)
		SetPermissions(userID, groupID string, permissions map[string]Permissions) (interface{}, error)
		getHost() *url.URL
	}

	gatekeeperClient struct {
		httpClient    *http.Client    // store a reference to the http client so we can reuse it
		hostGetter    disc.HostGetter // The getter that provides the host to talk to for the client
		tokenProvider TokenProvider   // An object that provides tokens for communicating with gatekeeper
	}

	gatekeeperClientBuilder struct {
		httpClient    *http.Client    // store a reference to the http client so we can reuse it
		hostGetter    disc.HostGetter // The getter that provides the host to talk to for the client
		tokenProvider TokenProvider   // An object that provides tokens for communicating with gatekeeper
	}

	Permissions map[string]interface{}
)

func NewGatekeeperClientBuilder() *gatekeeperClientBuilder {
	return &gatekeeperClientBuilder{}
}

func (b *gatekeeperClientBuilder) WithHttpClient(httpClient *http.Client) *gatekeeperClientBuilder {
	b.httpClient = httpClient
	return b
}

func (b *gatekeeperClientBuilder) WithHostGetter(hostGetter disc.HostGetter) *gatekeeperClientBuilder {
	b.hostGetter = hostGetter
	return b
}

func (b *gatekeeperClientBuilder) WithTokenProvider(tokenProvider TokenProvider) *gatekeeperClientBuilder {
	b.tokenProvider = tokenProvider
	return b
}

func (b *gatekeeperClientBuilder) Build() *gatekeeperClient {
	if b.hostGetter == nil {
		panic("gatekeeperClient requires a hostGetter to be set")
	}
	if b.tokenProvider == nil {
		panic("gatekeeperClient requires a tokenProvider to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &gatekeeperClient{
		httpClient:    b.httpClient,
		hostGetter:    b.hostGetter,
		tokenProvider: b.tokenProvider,
	}
}

func (client *gatekeeperClient) UserInGroup(userID, groupID string) (map[string]Permissions, error) {
	host := client.getHost()
	if host == nil {
		return nil, errors.New("No known gatekeeper hosts")
	}
	host.Path += fmt.Sprintf("access/%s/%s", groupID, userID)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

	log.Println(req)
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		retVal := make(map[string]Permissions)
		if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
			log.Println(err)
			return nil, &status.StatusError{status.NewStatus(500, "Unable to parse response.")}
		}
		return retVal, nil
	} else if res.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, &status.StatusError{status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}

}

func (client *gatekeeperClient) SetPermissions(userID, groupID string, permissions map[string]Permissions) (interface{}, error) {
	host := client.getHost()
	if host == nil {
		return nil, errors.New("No known gatekeeper hosts")
	}
	host.Path += fmt.Sprintf("access/%s/%s", groupID, userID)

	req, _ := http.NewRequest("POST", host.String(), bytes.NewBuffer(encodePermissions(permissions)))
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

	log.Println(req)
	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {

		var retVal interface{}
		if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
			errMsg := fmt.Sprintf("Unable to parse response: [%s]", err.Error())
			log.Println(errMsg)
			return nil, &status.StatusError{status.NewStatus(500, errMsg)}
		}
		return retVal, nil
	} else if res.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, &status.StatusError{status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}

}

func encodePermissions(permissions map[string]Permissions) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(permissions); err != nil {
		log.Println("Error encodePermissions ", err)
		return nil
	}
	return buf.Bytes()
}

func (client *gatekeeperClient) getHost() *url.URL {
	if hostArr := client.hostGetter.HostGet(); len(hostArr) > 0 {
		cpy := new(url.URL)
		*cpy = hostArr[0]
		return cpy
	} else {
		return nil
	}
}

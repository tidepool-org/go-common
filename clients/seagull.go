package clients

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path"

	"github.com/tidepool-org/go-common/clients/status"
)

type (
	Seagull interface {
		// Retrieves arbitrary collection information from metadata
		//
		// userID -- the Tidepool-assigned userId
		// hashName -- the name of what we are trying to get
		// token -- a server token or the user token
		GetPrivatePair(userID, hashName, token string) *PrivatePair
		// Retrieves arbitrary collection information from metadata
		//
		// userID -- the Tidepool-assigned userId
		// collectionName -- the collection being retrieved
		// token -- a server token or the user token
		// v - the interface to return the value in
		GetCollection(userID, collectionName, token string, v interface{}) error
	}

	seagullClient struct {
		httpClient *http.Client // store a reference to the http client so we can reuse it
		host       url.URL      // The host to talk to for the client
	}

	seagullClientBuilder struct {
		httpClient *http.Client
		host       *url.URL
		err        error
	}

	PrivatePair struct {
		ID    string
		Value string
	}
)

func NewSeagullClientBuilder() *seagullClientBuilder {
	return &seagullClientBuilder{}
}

func (b *seagullClientBuilder) WithHttpClient(httpClient *http.Client) *seagullClientBuilder {
	b.httpClient = httpClient
	return b
}

func (b *seagullClientBuilder) WithHost(host string) *seagullClientBuilder {
	h, err := url.Parse(host)
	if err != nil {
		b.err = err
	}
	b.host = h
	return b
}

func (b *seagullClientBuilder) Build() *seagullClient {
	if b.err != nil {
		panic(b.err)
	}
	if b.httpClient == nil {
		panic("seagullClient requires an httpClient to be set")
	}
	if b.host == nil {
		panic("seagullClient requires a host to be set")
	}
	return &seagullClient{
		httpClient: b.httpClient,
		host:       *b.host,
	}
}

func (client *seagullClient) GetPrivatePair(userID, hashName, token string) *PrivatePair {
	host := client.getHost()
	host.Path = path.Join(host.Path, userID, "private", hashName)

	req, err := http.NewRequest("GET", host.String(), nil)
	if err != nil {
		panic(err)
	}
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

func (client *seagullClient) GetCollection(userID, collectionName, token string, v interface{}) error {
	host := client.getHost()
	host.Path = path.Join(host.Path, userID, collectionName)

	req, err := http.NewRequest("GET", host.String(), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-tidepool-session-token", token)

	log.Println(req)
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

func (client *seagullClient) getHost() url.URL {
	return client.host
}

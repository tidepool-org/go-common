package clients

import (
	"encoding/json"
	"fmt"
	"github.com/tidepool-org/go-common/clients/disc"
	"log"
	"net/http"
	"net/url"
)

type (
	/*Seagull interface we export*/
	Seagull interface {
		GetPrivatePair(userID, hashName, token string) *PrivatePair
		GetCollection(userID, collectionName, token string) interface{}
		getHost() *url.URL
	}

	seagullClient struct {
		httpClient *http.Client    // store a reference to the http client so we can reuse it
		hostGetter disc.HostGetter // The getter that provides the host to talk to for the client
	}

	seagullClientBuilder struct {
		httpClient *http.Client
		hostGetter disc.HostGetter
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

func (b *seagullClientBuilder) WithHostGetter(hostGetter disc.HostGetter) *seagullClientBuilder {
	b.hostGetter = hostGetter
	return b
}

func (b *seagullClientBuilder) Build() *seagullClient {
	if b.httpClient == nil {
		panic("seagullClient requires an httpClient to be set")
	}
	if b.hostGetter == nil {
		panic("seagullClient requires a hostGetter to be set")
	}
	return &seagullClient{
		httpClient: b.httpClient,
		hostGetter: b.hostGetter,
	}
}

func (client *seagullClient) GetPrivatePair(userID, hashName, token string) *PrivatePair {
	host := client.getHost()
	if host == nil {
		return nil
	}
	host.Path += fmt.Sprintf("%s/private/%s", userID, hashName)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	log.Println(req)
	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Problem when looking up private pair for userID[%s]. %s", userID, err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Unknown response code[%s] from service[%s]", res.StatusCode, req.URL)
		return nil
	}

	var retVal PrivatePair
	if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
		log.Println("Error parsing JSON results", err)
		return nil
	}
	return &retVal
}

/*
 *	Retrieves arbitrary collection information from metadata
 *
 *  userID -- the Tidepool-assigned userId
 *  collectionName -- the collection being retrieved
 *  token -- a server token or the user token
 */
func (client *seagullClient) GetCollection(userID, collectionName, token string) interface{} {
	host := client.getHost()
	if host == nil {
		return nil
	}
	host.Path += fmt.Sprintf("/metadata/%s/%s", userID, collectionName)

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", token)

	log.Println(req)
	res, err := client.httpClient.Do(req)
	if err != nil {
		log.Printf("Problem when looking up collection for userID[%s]. %s", userID, err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK || res.StatusCode != http.StatusNoContent {
		log.Printf("Unknown response code[%s] from service[%s]", res.StatusCode, req.URL)
		return nil
	}

	var retVal interface{}
	if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
		log.Println("Error parsing JSON results", err)
		return nil
	}
	return &retVal
}

func (client *seagullClient) getHost() *url.URL {
	if hostArr := client.hostGetter.HostGet(); len(hostArr) > 0 {
		cpy := new(url.URL)
		*cpy = hostArr[0]
		return cpy
	} else {
		return nil
	}
}

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
	SourceArray []*Source
	//Inteface so that we can mock dataClient for tests
	Data interface {
		//userID  -- the Tidepool-assigned userID
		//
		// returns the Data Sources for the user
		ListSources(userID string) (SourceArray, error)
	}

	DataClient struct {
		httpClient    *http.Client    // store a reference to the http client so we can reuse it
		hostGetter    disc.HostGetter // The getter that provides the host to talk to for the client
		tokenProvider TokenProvider   // An object that provides tokens for communicating with data
	}

	dataClientBuilder struct {
		httpClient    *http.Client    // store a reference to the http client so we can reuse it
		hostGetter    disc.HostGetter // The getter that provides the host to talk to for the client
		tokenProvider TokenProvider   // An object that provides tokens for communicating with data
	}
)

type Source struct {
	ID                *string              `json:"id,omitempty" bson:"id,omitempty"`
	UserID            *string              `json:"userId,omitempty" bson:"userId,omitempty"`
	ProviderType      *string              `json:"providerType,omitempty" bson:"providerType,omitempty"`
	ProviderName      *string              `json:"providerName,omitempty" bson:"providerName,omitempty"`
	ProviderSessionID *string              `json:"providerSessionId,omitempty" bson:"providerSessionId,omitempty"`
	State             *string              `json:"state,omitempty" bson:"state,omitempty"`
	Error             *errors.Serializable `json:"error,omitempty" bson:"error,omitempty"`
	DataSetIDs        *[]string            `json:"dataSetIds,omitempty" bson:"dataSetIds,omitempty"`
	EarliestDataTime  *time.Time           `json:"earliestDataTime,omitempty" bson:"earliestDataTime,omitempty"`
	LatestDataTime    *time.Time           `json:"latestDataTime,omitempty" bson:"latestDataTime,omitempty"`
	LastImportTime    *time.Time           `json:"lastImportTime,omitempty" bson:"lastImportTime,omitempty"`
	CreatedTime       *time.Time           `json:"createdTime,omitempty" bson:"createdTime,omitempty"`
	ModifiedTime      *time.Time           `json:"modifiedTime,omitempty" bson:"modifiedTime,omitempty"`
	Revision          *int                 `json:"revision,omitempty" bson:"revision,omitempty"`
}

func NewDataClientBuilder() *dataClientBuilder {
	return &dataClientBuilder{}
}

func (b *dataClientBuilder) WithHttpClient(httpClient *http.Client) *dataClientBuilder {
	b.httpClient = httpClient
	return b
}

func (b *dataClientBuilder) WithHostGetter(hostGetter disc.HostGetter) *dataClientBuilder {
	b.hostGetter = hostGetter
	return b
}

func (b *dataClientBuilder) WithTokenProvider(tokenProvider TokenProvider) *dataClientBuilder {
	b.tokenProvider = tokenProvider
	return b
}

func (b *dataClientBuilder) Build() *DataClient {
	if b.hostGetter == nil {
		panic("dataClient requires a hostGetter to be set")
	}
	if b.tokenProvider == nil {
		panic("dataClient requires a tokenProvider to be set")
	}

	if b.httpClient == nil {
		b.httpClient = http.DefaultClient
	}

	return &DataClient{
		httpClient:    b.httpClient,
		hostGetter:    b.hostGetter,
		tokenProvider: b.tokenProvider,
	}
}

func (client *DataClient) ListSources(userID string) (SourceArray, error) {
	host := client.getHost()
	if host == nil {
		return nil, errors.New("No known data hosts")
	}
	host.Path = path.Join(host.Path, "/v1/users/", userID, "/data_sources")

	req, _ := http.NewRequest("GET", host.String(), nil)
	req.Header.Add("x-tidepool-session-token", client.tokenProvider.TokenProvide())

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		retVal := SourceArray{}
		if err := json.NewDecoder(res.Body).Decode(&retVal); err != nil {
			log.Println(err)
			return nil, &status.StatusError{status.NewStatus(500, "ListSources Unable to parse response.")}
		}
		return retVal, nil
	} else if res.StatusCode == 404 {
		return nil, nil
	} else {
		return nil, &status.StatusError{status.NewStatusf(res.StatusCode, "Unknown response code from service[%s]", req.URL)}
	}
}

func (client *DataClient) getHost() *url.URL {
	if hostArr := client.hostGetter.HostGet(); len(hostArr) > 0 {
		cpy := new(url.URL)
		*cpy = hostArr[0]
		return cpy
	} else {
		return nil
	}
}

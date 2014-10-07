package clients

import (
	"encoding/json"
	"net/url"

	"github.com/tidepool-org/go-common/clients/status"
)

type (
	seagullMock struct {
		fakeJSON string
		fakeUrl  *url.URL
	}
)

func NewSeagullMock(fakeJSON string, fakeUrl *url.URL) *seagullMock {
	return &seagullMock{
		fakeJSON: fakeJSON,
		fakeUrl:  fakeUrl,
	}

}

func (mock *seagullMock) GetPrivatePair(userID, hashName, token string) *PrivatePair {
	return &PrivatePair{ID: "mock.id.123", Value: "mock value"}
}

func (mock *seagullMock) GetCollection(userID, collectionName, token string, v interface{}) error {

	if mock.fakeJSON == "" {
		return &status.StatusError{status.NewStatus(500, "Unable to get collection.")}
	}

	json.Unmarshal([]byte(mock.fakeJSON), &v)

	return nil
}
func (mock *seagullMock) getHost() *url.URL {
	return mock.fakeUrl
}

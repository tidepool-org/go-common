package clients

import (
	"net/url"

	"github.com/tidepool-org/go-common/clients/status"
)

type (
	seagullMock struct {
		fakeCollection interface{}
		fakeUrl        *url.URL
	}
)

func NewSeagullMock(fakeCollection interface{}, fakeUrl *url.URL) *seagullMock {
	return &seagullMock{
		fakeCollection: fakeCollection,
		fakeUrl:        fakeUrl,
	}
}

func (mock *seagullMock) GetPrivatePair(userID, hashName, token string) *PrivatePair {
	return &PrivatePair{ID: "mock.id.123", Value: "mock value"}
}

func (mock *seagullMock) GetCollection(userID, collectionName, token string, v interface{}) error {

	if mock.fakeCollection == nil {
		return &status.StatusError{status.NewStatus(500, "Unable to get collection.")}
	}

	v = mock.fakeCollection

	return nil
}
func (mock *seagullMock) getHost() *url.URL {
	return mock.fakeUrl
}

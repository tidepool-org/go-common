package clients

import (
	"net/url"

	"github.com/tidepool-org/go-common/clients/status"
)

type (
	gatekeeperMock struct {
		expectedPermissons map[string]Permissions //what the mock will return for UserInGroup
		expectedHost       *url.URL               //what the mock will return for getHost
	}
)

func NewGatekeeperMock(permissonsToReturn map[string]Permissions, hostToReturn *url.URL) *gatekeeperMock {
	return &gatekeeperMock{
		expectedPermissons: permissonsToReturn,
		expectedHost:       hostToReturn,
	}
}

func (mock *gatekeeperMock) UserInGroup(userID, groupID string) (map[string]Permissions, error) {

	if mock.expectedPermissons == nil {
		return nil, &status.StatusError{status.NewStatus(500, "Unable to parse response.")}
	}
	return mock.expectedPermissons, nil
}

func (mock *gatekeeperMock) SetPermissions(userID, groupID string, permissions Permissions) (map[string]Permissions, error) {
	return mock.expectedPermissons, nil
}

func (mock *gatekeeperMock) getHost() *url.URL {
	return mock.expectedHost
}

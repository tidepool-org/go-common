package clients

import (
	"net/url"

	"github.com/tidepool-org/go-common/clients/status"
)

type (
	gatekeeperMock struct {
		ExpectedPermissons map[string]Permissions //what the mock will return for UserInGroup
		ExpectedHost       *url.URL               //what the mock will return for getHost
	}
)

func NewGatekeeperMock(permissonsToReturn map[string]Permissions, hostToReturn *url.URL) *gatekeeperMock {
	return &gatekeeperClient{
		ExpectedPermissons: permissonsToReturn,
		ExpectedHost:       hostToReturn,
	}
}

func (mock *gatekeeperMock) UserInGroup(userID, groupID string) (map[string]Permissions, error) {

	if mock.ExpectedPermissons == nil {
		return nil, &status.StatusError{status.NewStatus(500, "Unable to parse response.")}
	}
	return mock.ExpectedPermissons, nil
}

func (mock *gatekeeperMock) getHost() *url.URL {
	return mock.ExpectedHost
}

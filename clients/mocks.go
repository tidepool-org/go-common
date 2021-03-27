package clients

import (
	"context"
	"encoding/json"
)

type (
	GatekeeperMock struct {
		expectedPermissions Permissions
		expectedError       error
	}
	SeagullMock struct{}
)

//A mock of the Gatekeeper interface
func NewGatekeeperMock(expectedPermissions Permissions, expectedError error) *GatekeeperMock {
	return &GatekeeperMock{expectedPermissions, expectedError}
}

func (mock *GatekeeperMock) UserInGroup(ctx context.Context, userID, groupID string) (Permissions, error) {
	if mock.expectedPermissions != nil || mock.expectedError != nil {
		return mock.expectedPermissions, mock.expectedError
	} else {
		return Permissions{userID: Allowed}, nil
	}
}

func (mock *GatekeeperMock) UsersInGroup(ctx context.Context, groupID string) (UsersPermissions, error) {
	if mock.expectedPermissions != nil || mock.expectedError != nil {
		return UsersPermissions{groupID: mock.expectedPermissions}, mock.expectedError
	} else {
		return UsersPermissions{groupID: Permissions{groupID: Allowed}}, nil
	}
}

func (mock *GatekeeperMock) SetPermissions(ctx context.Context, userID, groupID string, permissions Permissions) (Permissions, error) {
	return Permissions{userID: Allowed}, nil
}

//A mock of the Seagull interface
func NewSeagullMock() *SeagullMock {
	return &SeagullMock{}
}

func (mock *SeagullMock) GetPrivatePair(ctx context.Context, userID, hashName, token string) *PrivatePair {
	return &PrivatePair{ID: "mock.id.123", Value: "mock value"}
}

func (mock *SeagullMock) GetCollection(ctx context.Context, userID, collectionName, token string, v interface{}) error {
	json.Unmarshal([]byte(`{"Something":"anit no thing", "patient": {"birthday": "2016-01-01"}}`), &v)
	return nil
}

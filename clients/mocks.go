package clients

import (
	"encoding/json"
)

type (
	GatekeeperMock struct {
		expectedPermissions Permissions
		expectedError       error
	}
	SeagullMock struct{}
	DataMock    struct {
		dataSources map[string]DataSourceArray
		dataSets    map[string]DataSetArray
	}
)

// A mock of the Gatekeeper interface
func NewGatekeeperMock(expectedPermissions Permissions, expectedError error) *GatekeeperMock {
	return &GatekeeperMock{expectedPermissions, expectedError}
}

func (mock *GatekeeperMock) UserInGroup(userID, groupID string) (Permissions, error) {
	if mock.expectedPermissions != nil || mock.expectedError != nil {
		return mock.expectedPermissions, mock.expectedError
	} else {
		return Permissions{userID: Allowed}, nil
	}
}

func (mock *GatekeeperMock) UsersInGroup(groupID string) (UsersPermissions, error) {
	if mock.expectedPermissions != nil || mock.expectedError != nil {
		return UsersPermissions{groupID: mock.expectedPermissions}, mock.expectedError
	} else {
		return UsersPermissions{groupID: Permissions{groupID: Allowed}}, nil
	}
}

func (mock *GatekeeperMock) GroupsForUser(userID string) (UsersPermissions, error) {
	if mock.expectedPermissions != nil || mock.expectedError != nil {
		return UsersPermissions{userID: mock.expectedPermissions}, mock.expectedError
	} else {
		return UsersPermissions{userID: Permissions{userID: Allowed}}, nil
	}
}

func (mock *GatekeeperMock) SetPermissions(userID, groupID string, permissions Permissions) (Permissions, error) {
	return Permissions{userID: Allowed}, nil
}

// A mock of the Seagull interface
func NewSeagullMock() *SeagullMock {
	return &SeagullMock{}
}

func (mock *SeagullMock) GetPrivatePair(userID, hashName, token string) *PrivatePair {
	return &PrivatePair{ID: "mock.id.123", Value: "mock value"}
}

func (mock *SeagullMock) GetCollection(userID, collectionName, token string, v interface{}) error {
	json.Unmarshal([]byte(`{"Something":"anit no thing", "patient": {"birthday": "2016-01-01"}}`), &v)
	return nil
}

func (mock *SeagullMock) UpdateCollection(userID, collectionName, token string, v interface{}) error {
	return nil
}

func NewDataMock(userDataSources map[string]DataSourceArray, userDataSets map[string]DataSetArray) *DataMock {
	if userDataSources == nil {
		userDataSources = map[string]DataSourceArray{}
	}
	if userDataSets == nil {
		userDataSets = map[string]DataSetArray{}
	}
	return &DataMock{
		dataSources: userDataSources,
		dataSets:    userDataSets,
	}
}

func (m *DataMock) ListSources(userID string) (DataSourceArray, error) {
	return m.dataSources[userID], nil
}

func (m *DataMock) ListSourcesPagination(userID string, p Pagination) (DataSourceArray, error) {
	return m.dataSources[userID], nil
}

func (m *DataMock) ListSetsPagination(userID string, p Pagination) (DataSetArray, error) {
	return m.dataSets[userID], nil
}

func (m *DataMock) HasAnyData(userID string) (hasData bool, err error) {
	return len(m.dataSets[userID]) > 0, nil
}

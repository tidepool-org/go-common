package shoreline

import (
	"log"
)

type ShorelineMockClient struct {
	ServerSecret string
}

func NewMock(secret string) *ShorelineMockClient {
	return &ShorelineMockClient{ServerSecret: secret}
}

func (client *ShorelineMockClient) Start() error {
	log.Println("Started mock shoreline client")
	return nil
}

func (client *ShorelineMockClient) Close() {
	log.Println("Close mock shoreline client")
}

func (client *ShorelineMockClient) Login(username, password string) (*UserData, string, error) {
	return &UserData{UserID: "123.456.789", Username: username, Emails: []string{username}}, client.ServerSecret, nil
}

func (client *ShorelineMockClient) Signup(username, password, email string) (*UserData, error) {
	return &UserData{UserID: "123.xxx.456", Username: username, Emails: []string{email}}, nil
}

func (client *ShorelineMockClient) CheckToken(token string) *TokenData {
	return &TokenData{UserID: "987.654.321", IsServer: true}
}

func (client *ShorelineMockClient) CheckTokenForScopes(requiredScopes, token string) *TokenData {
	return &TokenData{UserID: "987.654.321", IsServer: true}
}

func (client *ShorelineMockClient) SecretProvide() string {
	return client.ServerSecret
}

func (client *ShorelineMockClient) GetUser(userID, token string) (*UserData, error) {
	if userID == "NotFound" {
		return nil, nil
	} else if userID == "WithoutPassword" {
		return &UserData{UserID: userID, Username: "From Mock", Emails: []string{userID}, PasswordExists: false}, nil
	} else {
		return &UserData{UserID: userID, Username: "From Mock", Emails: []string{userID}, PasswordExists: true}, nil
	}
}

func (client *ShorelineMockClient) UpdateUser(userID string, userUpdate UserUpdate, token string) error {
	return nil
}

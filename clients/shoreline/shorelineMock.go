package shoreline

import (
	"context"
	"log"
)

type ShorelineMockClient struct {
	ServerToken string
}

func NewMock(token string) *ShorelineMockClient {
	return &ShorelineMockClient{ServerToken: token}
}

func (client *ShorelineMockClient) Start(ctx context.Context) error {
	log.Println("Started mock shoreline client")
	return nil
}

func (client *ShorelineMockClient) Close(ctx context.Context) {
	log.Println("Close mock shoreline client")
}

func (client *ShorelineMockClient) Login(ctx context.Context, username, password string) (*UserData, string, error) {
	return &UserData{UserID: "123.456.789", Username: username, Emails: []string{username}}, client.ServerToken, nil
}

func (client *ShorelineMockClient) Signup(ctx context.Context, username, password, email string) (*UserData, error) {
	return &UserData{UserID: "123.xxx.456", Username: username, Emails: []string{email}}, nil
}

func (client *ShorelineMockClient) CheckToken(ctx context.Context, token string) *TokenData {
	return &TokenData{UserID: "987.654.321", IsServer: true}
}

func (client *ShorelineMockClient) TokenProvide(ctx context.Context) string {
	return client.ServerToken
}

func (client *ShorelineMockClient) GetUser(ctx context.Context, userID, token string) (*UserData, error) {
	if userID == "NotFound" {
		return nil, nil
	} else if userID == "WithoutPassword" {
		return &UserData{UserID: userID, Username: "From Mock", Emails: []string{userID}, PasswordExists: false}, nil
	} else {
		return &UserData{UserID: userID, Username: "From Mock", Emails: []string{userID}, PasswordExists: true}, nil
	}
}

func (client *ShorelineMockClient) UpdateUser(ctx context.Context, userID string, userUpdate UserUpdate, token string) error {
	return nil
}

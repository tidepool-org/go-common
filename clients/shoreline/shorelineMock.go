package shoreline

import (
	"log"
	"net/url"
)

type ShorelineMockClient struct {
	serverToken string
}

/*
 	Start() error
	Close()
	serverLogin() error
	Login(username, password string) (*UserData, string, error)
	CheckToken(token string) *TokenData
	TokenProvide() string
	getHost() *url.URL
*/

func (client *ShorelineMockClient) Start() error {
	log.Println("Started mock shoreline client")
	return nil
}

func (client *ShorelineMockClient) Close() {
	log.Println("Close mock shoreline client")
}

func (client *ShorelineMockClient) serverLogin() error {
	client.serverToken = "a.mock.token"
	return nil
}

func (client *ShorelineMockClient) Login(username, password string) (*UserData, string, error) {
	return &UserData{UserID: "123.456.789", UserName: username, Emails: []string{username}}, client.serverToken, nil
}

func (client *ShorelineMockClient) CheckToken(token string) *TokenData {
	return &TokenData{UserID: "987.654.321", IsServer: true}
}

func (client *ShorelineMockClient) TokenProvide() string {
	return client.serverToken
}

func (client *ShorelineMockClient) getHost() *url.URL {
	return &url.URL{}
}

package highwater

import (
	"log"
)

type HighwaterMockClient struct{}

func NewMock() *HighwaterMockClient {
	return &HighwaterMockClient{}
}

func (client *HighwaterMockClient) postServer(eventName, token string, params map[string]string) {

	if eventName == "" || token == "" || len(params) <= 0 {
		log.Panicf("missing required eventName[%s] token[%s] params[%v]", eventName, token, params)
	}

	return
}

func (client *HighwaterMockClient) postThisUser(eventName, token string, params map[string]string) {

	if eventName == "" || token == "" || len(params) <= 0 {
		log.Panicf("missing required eventName[%s] token[%s] params[%v]", eventName, token, params)
	}

	return
}

func (client *HighwaterMockClient) postWithUser(userId, eventName, token string, params map[string]string) {
	if userId == "" || eventName == "" || token == "" || len(params) <= 0 {
		log.Panicf("missing required userId[%s] eventName[%s] token[%s] params[%v]", userId, eventName, token, params)
	}

	return
}

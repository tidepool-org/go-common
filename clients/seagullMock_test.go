package clients

import (
	"net/url"
	"testing"
)

//ensure the mock will actually be usefull
func TestMock(t *testing.T) {

	const TOKEN_MOCK = "this is a token"
	mockUrl, _ := url.Parse("http://something.org/search?q=yay")

	seagullClient := NewSeagullMock(`{"Something":"anit no thing"}`, mockUrl)

	var fake struct{ Something string }
	seagullClient.GetCollection("123.456", "stuff", TOKEN_MOCK, &fake)

	if fake.Something != "anit no thing" {
		t.Error("Should have given mocked collection")
	}

	if pp := seagullClient.GetPrivatePair("123.456", "Stuff", TOKEN_MOCK); pp == nil {
		t.Error("Should give us mocked private pair")
	}

	if host := seagullClient.getHost(); host == nil {
		t.Error("Should give us a mocked host")
	}

}

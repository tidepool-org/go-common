package clients

import (
	"net/url"
	"testing"
)

//ensure the mock will actually be usefull
func TestMock(t *testing.T) {

	const TOKEN_MOCK = "this is a token"
	mockUrl, _ := url.Parse("http://something.org/search?q=yay")
	type Fake struct{ Something string }

	seagullClient := NewSeagullMock(&Fake{Something: "anit no thing"}, mockUrl)

	fd := &Fake{}
	seagullClient.GetCollection("123.456", "stuff", TOKEN_MOCK, &fd)

	t.Logf("coll [%v]", fd)

	//if fd.Something == "" {
	//	t.Error("Should have given mocked collection")
	//}

	if pp := seagullClient.GetPrivatePair("123.456", "Stuff", TOKEN_MOCK); pp == nil {
		t.Error("Should give us mocked private pair")
	}

	if host := seagullClient.getHost(); host == nil {
		t.Error("Should give us a mocked host")
	}

}

package shoreline

import (
	"fmt"
	"github.com/tidepool-org/go-common/clients/disc"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const token = "this is a token"

func TestMock(t *testing.T) {

	client := &ShorelineMockClient{serverToken: token}

	if err := client.Start(); err != nil {
		t.Errorf("Failed start with error[%v]", err)
	}

	if tok := client.TokenProvide(); tok != token {
		t.Errorf("Unexpected token[%s]", tok)
	}

	client.Close()

	if usr, token, err := client.Login("billy", "howdy"); err != nil {
		t.Errorf("Failed start with error[%v]", err)
	} else {
		if usr == nil {
			t.Error("Should give us a fake usr details")
		}
		if token == "" {
			t.Error("Should give us a fake token")
		}
	}

	if host := client.getHost(); host == nil {
		t.Error("Should give us a fake host")
	}

	if err := client.serverLogin(); err != nil {
		t.Error("Should not return err")
	}

	if svrToken := client.serverToken(); svrToken == nil {
		t.Error("Should give us a token")
	}

}

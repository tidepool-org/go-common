package clients

import (
	"net/url"
	"testing"
)

const USERID, GROUPID = "123user", "456group"

func TestGatekeeperMock(t *testing.T) {

	mockUrl, _ := url.Parse("http://something.org/search?q=yay")
	permissonsToReturn := make(map[string]Permissions)
	p := make(Permissions)

	p["userid"] = USERID

	permissonsToReturn["root"] = p

	gatekeeperMock := NewGatekeeperMock(permissonsToReturn, mockUrl)

	if perms, err := gatekeeperMock.UserInGroup(USERID, GROUPID); err != nil {
		t.Fatal("No error should be returned")
	} else {
		if perms == nil {
			t.Fatalf("Perms where [%v] but expected [%v]", perms, permissonsToReturn)
		}
		t.Logf("Perms where [%v] given [%v]", perms, permissonsToReturn)
	}

	if host := gatekeeperMock.getHost(); host != mockUrl {
		t.Fatalf("Host was [%v] but expected [%v]", host, mockUrl)
	}

	if perms, err := gatekeeperMock.SetPermissions(USERID, GROUPID, permissonsToReturn["root"]); err != nil {
		t.Fatal("No error should be returned")
	} else {
		if perms == nil {
			t.Fatalf("Perms where [%v] but expected [%v]", perms, permissonsToReturn["root"])
		}
		t.Logf("Perms where [%v] given [%v]", perms, permissonsToReturn)
	}

}

func TestGatekeeperMock_WhenNil(t *testing.T) {

	gatekeeperMock := NewGatekeeperMock(nil, nil)

	if perms, err := gatekeeperMock.UserInGroup(USERID, GROUPID); err == nil {
		t.Fatal("There should have been an error returned")
	} else {
		if perms != nil {
			t.Fatalf("Perms where [%v] but expected none", perms)
		}
	}

	if host := gatekeeperMock.getHost(); host != nil {
		t.Fatalf("Host was [%v] but expected none", host)
	}

}

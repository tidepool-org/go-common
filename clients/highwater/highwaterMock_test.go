package highwater

import (
	"context"
	"testing"
)

const (
	EVENT_NAME = "testing"
	USERID     = "123-456-cc-2"
	TOKEN      = "a.fake.token.for.this"
)

func TestMock(t *testing.T) {

	p := make(map[string]string)

	p["one"] = "two"
	p["buckle"] = "my"
	p["shoe"] = "three ..."

	client := NewMock()

	client.PostServer(context.Background(), EVENT_NAME, TOKEN, p)

	client.PostThisUser(context.Background(), EVENT_NAME, TOKEN, p)

	client.PostWithUser(context.Background(), USERID, EVENT_NAME, TOKEN, p)
}

//log.Panic is called when not all required args are passed.
//This test fails if the panic and subseqent recover are not called
func TestMock_Fails(t *testing.T) {

	defer func() {
		if r := recover(); r == nil {
			t.Error("should have paniced")
		}
	}()

	p := make(map[string]string)

	p["one"] = "two"
	p["buckle"] = "my"
	p["shoe"] = "three ..."

	client := NewMock()

	client.PostServer(context.Background(), "", TOKEN, p)

	client.PostThisUser(context.Background(), EVENT_NAME, "", p)

	client.PostWithUser(context.Background(), "", EVENT_NAME, TOKEN, p)

	client.PostWithUser(context.Background(), "", EVENT_NAME, TOKEN, nil)
}

// package clients is a set of structs and methods for client libraries that interact with the various
// services in the tidepool platform
package clients

import "context"

type TokenProvider interface {
	TokenProvide(ctx context.Context) string
}

type TokenProviderFunc func() string

func (t TokenProviderFunc) TokenProvide(ctx context.Context) string {
	return t()
}

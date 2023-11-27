package hooks

import (
	"context"
	"errors"
	"net/url"
	"os"

	mh "github.com/mochi-mqtt/hooks/auth/http"

	"google.golang.org/api/idtoken"
)

func NewMochiCloudHooksHTTPAuthConfig(ctx context.Context) (*mh.Options, error) {
	endpoint, found := os.LookupEnv("HTTP_AUTH_BACKEND_API_ENDPOINT")
	if !found {
		return nil, errors.New("HTTP_AUTH_BACKEND_API_ENDPOINT not found")
	}

	client, err := idtoken.NewClient(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return &mh.Options{
		ACLHost:                  stringToURL(endpoint + "/api/v1/auth/acls"),
		SuperUserHost:            stringToURL(endpoint + "/api/v1/clients/superuser"),
		ClientAuthenticationHost: stringToURL(endpoint + "/api/v1/auth/clients"),
		RoundTripper:             client.Transport,
	}, nil
}

func stringToURL(s string) *url.URL {
	parsedURL, _ := url.Parse(s)
	return parsedURL
}

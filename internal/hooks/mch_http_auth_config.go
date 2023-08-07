package hooks

import (
	"context"
	"errors"
	"os"

	mch "github.com/dgduncan/mochi-cloud-hooks"
	"google.golang.org/api/idtoken"
)

func NewMochiCloudHooksHTTPAuthConfig(ctx context.Context) (*mch.HTTPAuthHookConfig, error) {
	endpoint, found := os.LookupEnv("HTTP_AUTH_BACKEND_API_ENDPOINT")
	if !found {
		return nil, errors.New("HTTP_AUTH_BACKEND_API_ENDPOINT not found")
	}

	client, err := idtoken.NewClient(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return &mch.HTTPAuthHookConfig{
		ACLHost:                  endpoint + "/api/v1/auth/acls",
		SuperUserHost:            endpoint + "/api/v1/clients/superuser",
		ClientAuthenticationHost: endpoint + "/api/v1/auth/clients",
		RoundTripper:             client.Transport,
	}, nil
}

package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/upperz-llc/go-broker/pkg/domain"
	"google.golang.org/api/idtoken"
)

type MochiBrokerAPIHTTP struct {
	HTTPClient *http.Client
	Endpoint   string
}

func (mhttp *MochiBrokerAPIHTTP) GetClientConnectionStatus(ctx context.Context, clientID string) (domain.ConnectionStatusGETResponse, error) {
	u, _ := url.Parse(mhttp.Endpoint)
	u.Path = path.Join(u.Path, "api", "v1", "client", clientID, "connection_status")

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)

	resp, err := mhttp.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return domain.ConnectionStatusGETResponse{}, err
	}
	defer resp.Body.Close()

	var cs domain.ConnectionStatusGETResponse
	json.NewDecoder(resp.Body).Decode(&cs)

	if resp.StatusCode == http.StatusNotFound {
		return domain.ConnectionStatusGETResponse{Connected: false}, nil
	}

	if resp.StatusCode != 200 {
		return domain.ConnectionStatusGETResponse{}, fmt.Errorf("non-200 and non-404 status code returned: %d", resp.StatusCode)
	}

	return cs, nil
}

func NewMochiBrokerAPIHTTP(ctx context.Context, hc *http.Client) (*MochiBrokerAPIHTTP, error) {
	endpoint, found := os.LookupEnv("MOCHI_BROKER_API_HTTP_ENDPOINT")
	if !found {
		return nil, errors.New("MOCHI_BROKER_API_HTTP_ENDPOINT not found")
	}
	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	var httpclient *http.Client
	if hc != nil {
		httpclient = hc

	} else {
		client, err := idtoken.NewClient(ctx, endpoint)
		if err != nil {
			return nil, err
		}
		httpclient = client

	}

	return &MochiBrokerAPIHTTP{
		HTTPClient: httpclient,
		Endpoint:   endpoint,
	}, nil
}

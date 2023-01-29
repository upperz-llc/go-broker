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

func (mhttp *MochiBrokerAPIHTTP) GetClientConnectionStatus(ctx context.Context, clientID string) (domain.OnlineStatusGETResponse, error) {
	u, _ := url.Parse(mhttp.Endpoint)
	u.Path = path.Join(u.Path, "api", "v1", "client", clientID, "online")

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)

	resp, err := mhttp.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return domain.OnlineStatusGETResponse{}, err
	}
	defer resp.Body.Close()

	var cs domain.OnlineStatusGETResponse
	json.NewDecoder(resp.Body).Decode(&cs)

	if resp.StatusCode == http.StatusNotFound {
		return domain.OnlineStatusGETResponse{Connected: false}, nil
	}

	if resp.StatusCode != 200 {
		return domain.OnlineStatusGETResponse{}, fmt.Errorf("non-200 and non-404 status code returned: %d", resp.StatusCode)
	}

	return cs, nil
}

// func (mhttp *MosquittoHTTP) SendConfigToDevice(ctx context.Context, deviceID string, config domain.DeviceConfig) error {
// 	u, _ := url.Parse(mhttp.Endpoint)
// 	u.Path = path.Join(u.Path, "device", deviceID, "config")

// 	body, _ := json.Marshal(config)

// 	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(body))

// 	resp, err := mhttp.HTTPClient.Do(req.WithContext(ctx))
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	log.Println(resp.StatusCode)

// 	return nil
// }

// func (mhttp *MosquittoHTTP) GetDeviceLWT(ctx context.Context, deviceID string) (domain.MosquittoAppLWTResponse, error) {
// 	u, _ := url.Parse(mhttp.Endpoint)
// 	u.Path = path.Join(u.Path, "device", deviceID, "lwt")

// 	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), http.NoBody)

// 	resp, err := mhttp.HTTPClient.Do(req.WithContext(ctx))
// 	if err != nil {
// 		return domain.MosquittoAppLWTResponse{}, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return domain.MosquittoAppLWTResponse{}, fmt.Errorf("non-200 status code returned: %d", resp.StatusCode)
// 	}

// 	var lwtResponse domain.MosquittoAppLWTResponse
// 	json.NewDecoder(resp.Body).Decode(&lwtResponse)

// 	return lwtResponse, nil
// }

func NewMochiBrokerAPIHTTP(ctx context.Context) (*MochiBrokerAPIHTTP, error) {
	endpoint, found := os.LookupEnv("MOCHI_BROKER_API_HTTP_ENDPOINT")
	if !found {
		return nil, errors.New("MOCHI_BROKER_API_HTTP_ENDPOINT not found")
	}
	// client is a http.Client that automatically adds an "Authorization" header
	// to any requests made.
	client, err := idtoken.NewClient(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	return &MochiBrokerAPIHTTP{
		HTTPClient: client,
		Endpoint:   endpoint,
	}, nil
}

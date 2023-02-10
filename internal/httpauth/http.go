package httpauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path"

	"google.golang.org/api/idtoken"
)

type HTTPAuthBackendHTTPIface interface {
	CheckClientACLs(ctx context.Context, payload ACLCheckPOST) (bool, error)
	CheckClientAuth(ctx context.Context, payload ClientCheckPOST) (bool, error)
	CheckSuperUserAuth(ctx context.Context, payload SuperuserCheckPOST) (bool, error)
}

type HTTPAuthBackendHTTP struct {
	HTTPClient *http.Client
	Endpoint   string
}

func (mhttp *HTTPAuthBackendHTTP) CheckSuperUserAuth(ctx context.Context, payload SuperuserCheckPOST) (bool, error) {
	u, _ := url.Parse(mhttp.Endpoint)
	u.Path = path.Join(u.Path, "api", "v1", "clients", "superuser")

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))

	resp, err := mhttp.HTTPClient.Do(req.WithContext(ctx))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func (mhttp *HTTPAuthBackendHTTP) CheckClientAuth(ctx context.Context, payload ClientCheckPOST) (bool, error) {
	u, _ := url.Parse(mhttp.Endpoint)
	u.Path = path.Join(u.Path, "api", "v1", "clients")

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))

	resp, err := mhttp.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func (mhttp *HTTPAuthBackendHTTP) CheckClientACLs(ctx context.Context, payload ACLCheckPOST) (bool, error) {
	u, _ := url.Parse(mhttp.Endpoint)
	u.Path = path.Join(u.Path, "api", "v1", "acls")

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(requestBody))

	resp, err := mhttp.HTTPClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func NewMochiBrokerAPIHTTP(ctx context.Context, hc *http.Client) (*HTTPAuthBackendHTTP, error) {
	endpoint, found := os.LookupEnv("HTTP_AUTH_BACKEND_API_ENDPOINT")
	if !found {
		return nil, errors.New("HTTP_AUTH_BACKEND_API_ENDPOINT not found")
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

	return &HTTPAuthBackendHTTP{
		HTTPClient: httpclient,
		Endpoint:   endpoint,
	}, nil
}

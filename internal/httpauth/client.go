package httpauth

import (
	"context"
	"net/http"
)

type SuperuserCheckPOST struct {
	Username string `json:"username"`
}

type ClientCheckPOST struct {
	ClientID string `json:"clientid"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type ACLCheckPOST struct {
	Username string `json:"username"`
	ClientID string `json:"clientid"`
	Topic    string `json:"topic"`
	ACC      string `json:"acc"`
}

type HTTPAuthBackendIface interface {
	CheckClientACLs(ctx context.Context, clientID, username, topic, acc string)
	CheckClientAuth(ctx context.Context, clientID, username, password string)
	CheckSuperUserAuth(ctx context.Context, username string)
}

type HTTPAuthBackendClient struct {
	HTTP HTTPAuthBackendHTTPIface
}

func (mbac *HTTPAuthBackendClient) CheckClientACLs(ctx context.Context, clientID, username, topic, acc string) (bool, error) {
	payload := ACLCheckPOST{
		ClientID: clientID,
		Username: username,
		Topic:    topic,
		ACC:      acc,
	}

	return mbac.HTTP.CheckClientACLs(ctx, payload)
}

func (mbac *HTTPAuthBackendClient) CheckClientAuth(ctx context.Context, clientID, username, password string) (bool, error) {
	payload := ClientCheckPOST{
		ClientID: clientID,
		Password: password,
		Username: username,
	}

	return mbac.HTTP.CheckClientAuth(ctx, payload)
}

func (mbac *HTTPAuthBackendClient) CheckSuperUserAuth(ctx context.Context, username string) (bool, error) {
	payload := SuperuserCheckPOST{
		Username: username,
	}

	return mbac.HTTP.CheckSuperUserAuth(ctx, payload)
}

func NewClient(ctx context.Context, hc *http.Client) (*HTTPAuthBackendClient, error) {
	mHTTP, err := NewMochiBrokerAPIHTTP(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &HTTPAuthBackendClient{
		HTTP: mHTTP,
	}, nil
}

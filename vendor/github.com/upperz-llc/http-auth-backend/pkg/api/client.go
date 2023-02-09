package api

import (
	"context"

	"github.com/upperz-llc/http-auth-backend/internal/domain"
)

type HTTPAuthBackendIface interface {
	CheckClientACLs(ctx context.Context, clientID, username, topic, acc string)
	CheckClientAuth(ctx context.Context, clientID, username, password string)
	CheckSuperUserAuth(ctx context.Context, username string)
}

type HTTPAuthBackendClient struct {
	HTTP HTTPAuthBackendHTTPIface
}

func (mbac *HTTPAuthBackendClient) CheckClientACLs(ctx context.Context, clientID, username, topic, acc string) (bool, error) {
	payload := domain.ACLCheckPOST{
		ClientID: clientID,
		Username: username,
		Topic:    topic,
		ACC:      acc,
	}

	return mbac.HTTP.CheckClientACLs(ctx, payload)
}

func (mbac *HTTPAuthBackendClient) CheckClientAuth(ctx context.Context, clientID, username, password string) (bool, error) {
	payload := domain.ClientCheckPOST{
		ClientID: clientID,
		Password: password,
		Username: username,
	}

	return mbac.HTTP.CheckClientAuth(ctx, payload)
}

func (mbac *HTTPAuthBackendClient) CheckSuperUserAuth(ctx context.Context, username string) (bool, error) {
	payload := domain.SuperuserCheckPOST{
		Username: username,
	}

	return mbac.HTTP.CheckSuperUserAuth(ctx, payload)
}

func NewClient(ctx context.Context) (*HTTPAuthBackendClient, error) {
	mHTTP, err := NewMochiBrokerAPIHTTP(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &HTTPAuthBackendClient{
		HTTP: mHTTP,
	}, nil
}

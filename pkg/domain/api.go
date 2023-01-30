package domain

import "context"

type BrokerAPIIface interface {
	GetClientConnectionStatus(ctx context.Context, clientID string) (bool, error)
}

type BrokerAPIHTTPIface interface {
	GetClientConnectionStatus(ctx context.Context, clientID string) (ConnectionStatusGETResponse, error)
}

type ConnectionStatusGETResponse struct {
	Connected bool `json:"connected"`
}

package domain

import "context"

type BrokerAPIIface interface {
	GetClientConnectionStatus(ctx context.Context, clientID string) (OnlineStatusGETResponse, error)
}

type OnlineStatusGETResponse struct {
	Connected bool `json:"connected"`
}

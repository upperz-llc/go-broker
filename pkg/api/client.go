package broker

import (
	"context"

	"github.com/upperz-llc/go-broker/pkg/domain"
)

type MochiBrokerAPIClient struct {
	HTTP domain.BrokerAPIIface
}

func (mbac *MochiBrokerAPIClient) GetClientConnectionStatus(ctx context.Context, clientID string) (bool, error) {
	onlinestatus, err := mbac.HTTP.GetClientConnectionStatus(ctx, clientID)
	if err != nil {
		return false, err
	}

	return onlinestatus.Connected, nil
}

func NewClient(ctx context.Context) (*MochiBrokerAPIClient, error) {
	mHTTP, err := NewMochiBrokerAPIHTTP(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &MochiBrokerAPIClient{
		HTTP: mHTTP,
	}, nil

}

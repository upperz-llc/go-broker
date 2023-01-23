package ps

import (
	"context"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetServiceFromServiceNameService(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bps := BrokerPubSub{}

	tests := []struct {
		name      string
		mocks     func(ctx context.Context)
		expectErr bool
	}{
		{
			name: "Success - Golden Path",
			mocks: func(ctx context.Context) {
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if tt.mocks != nil {
				tt.mocks(ctx)
			}

			err := bps.Publish(&pubsub.Topic{}, "")
			require.Nil(t, err)
		})
	}
}

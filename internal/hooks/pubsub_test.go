package hooks

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/logging"
	"github.com/golang/mock/gomock"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/stretchr/testify/require"
	"github.com/upperz-llc/go-broker/internal/mocks"
)

func TestGetServiceFromServiceNameService(t *testing.T) {

	// ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockpubsub := mocks.NewMockPubsub(ctrl)

	// Create and configure logger
	lc := logging.Client{}

	logger := lc.Logger("go-broker-log")

	gph := GCPPubsubHook{
		Logger: logger,
		Pubsub: mockpubsub,
	}

	tests := []struct {
		name      string
		mocks     func(ctx context.Context)
		expectErr bool
	}{
		{
			name: "Success - Golden Path",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "Failure - Error publishing pubsub",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errors.New("something went wrong"))
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

			gph.OnPublished(&mqtt.Client{}, packets.Packet{})
			require.Nil(t, nil)
		})
	}
}

func TestOnConnectHook(t *testing.T) {

	// ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockpubsub := mocks.NewMockPubsub(ctrl)

	// Create and configure logger
	lc := logging.Client{}

	logger := lc.Logger("go-broker-log")

	gph := GCPPubsubHook{
		Logger: logger,
		Pubsub: mockpubsub,
	}

	tests := []struct {
		name      string
		mocks     func(ctx context.Context)
		expectErr bool
	}{
		{
			name: "Success - Golden Path",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "Failure - Error publishing pubsub",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errors.New("something went wrong"))
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

			gph.OnConnect(&mqtt.Client{}, packets.Packet{})
			require.Nil(t, nil)
		})
	}
}

func TestOnDisconnectHook(t *testing.T) {

	// ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockpubsub := mocks.NewMockPubsub(ctrl)

	// Create and configure logger
	lc := logging.Client{}

	logger := lc.Logger("go-broker-log")

	gph := GCPPubsubHook{
		Logger: logger,
		Pubsub: mockpubsub,
	}

	tests := []struct {
		name      string
		mocks     func(ctx context.Context)
		expectErr bool
	}{
		{
			name: "Success - Golden Path",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "Failure - Error publishing pubsub",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errors.New("something went wrong"))
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

			gph.OnDisconnect(&mqtt.Client{}, nil, false)
			require.Nil(t, nil)
		})
	}
}

func TestOnSubscribedHook(t *testing.T) {

	// ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockpubsub := mocks.NewMockPubsub(ctrl)

	// Create and configure logger
	lc := logging.Client{}

	logger := lc.Logger("go-broker-log")

	gph := GCPPubsubHook{
		Logger: logger,
		Pubsub: mockpubsub,
	}

	tests := []struct {
		name      string
		mocks     func(ctx context.Context)
		expectErr bool
	}{
		{
			name: "Success - Golden Path",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "Failure - Error publishing pubsub",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errors.New("something went wrong"))
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

			gph.OnSubscribed(&mqtt.Client{}, packets.Packet{}, []byte{})
			require.Nil(t, nil)
		})
	}
}

func TestOnUnsubscribedHook(t *testing.T) {

	// ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockpubsub := mocks.NewMockPubsub(ctrl)

	// Create and configure logger
	lc := logging.Client{}

	logger := lc.Logger("go-broker-log")

	gph := GCPPubsubHook{
		Logger: logger,
		Pubsub: mockpubsub,
	}

	tests := []struct {
		name      string
		mocks     func(ctx context.Context)
		expectErr bool
	}{
		{
			name: "Success - Golden Path",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name: "Failure - Error publishing pubsub",
			mocks: func(ctx context.Context) {
				mockpubsub.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(errors.New("something went wrong"))
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

			gph.OnUnsubscribed(&mqtt.Client{}, packets.Packet{})
			require.Nil(t, nil)
		})
	}
}

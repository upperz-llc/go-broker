package hooks

import (
	"context"

	mch "github.com/dgduncan/mochi-cloud-hooks"
)

func NewMochiCloudHooksSecretManagerConfig(ctx context.Context) (*mch.SecretManagerHookConfig, error) {
	return &mch.SecretManagerHookConfig{
		Names: []string{"projects/481474188273/secrets/BROKER_ADMIN/versions/latest"},
	}, nil
}

package admin

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func getAdminCredentials(ctx context.Context) (string, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
		Name: "projects/freezer-monitor-dev-e7d4c/secrets/BROKER_ADMIN",
	})
	if err != nil {

		return "", err
	}

	return secret.String(), nil
}

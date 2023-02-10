package admin

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func getAdminCredentials(ctx context.Context) (string, error) {
	fmt.Println("getAdminCredentials")
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create secretmanager client: %v", err)
	}
	defer client.Close()

	// secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
	// 	Name: "projects/481474188273/secrets/BROKER_ADMIN",
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(secret.String())
	// secret.

	resp, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/481474188273/secrets/BROKER_ADMIN/versions/latest",
	})
	if err != nil {
		return "", err
	}

	return string(resp.Payload.Data), nil
}

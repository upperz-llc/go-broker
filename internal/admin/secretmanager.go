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

	resp, err := client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/481474188273/secrets/BROKER_ADMIN/versions/latest",
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(string(resp.Payload.String()))
	fmt.Println(string(resp.Payload.Data))

	resp, err = client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/freezer-monitor-dev-e7d4c/secrets/BROKER_ADMIN/versions/latest",
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(string(resp.Payload.String()))
	fmt.Println(string(resp.Payload.Data))

	resp, err = client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/freezer-monitor-dev-e7d4c/secrets/BROKER_ADMIN/versions/latest",
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	fmt.Println(string(resp.Payload.String()))
	fmt.Println(string(resp.Payload.Data))

	// secret, err := client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
	// 	Name: "projects/481474188273/secrets/BROKER_ADMIN",
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return "", err
	// }
	// fmt.Println(secret.String())

	// secret, err = client.GetSecret(ctx, &secretmanagerpb.GetSecretRequest{
	// 	Name: "BROKER_ADMIN",
	// })
	// if err != nil {
	// 	fmt.Println(err)
	// 	return "", err
	// }

	return "", nil
}

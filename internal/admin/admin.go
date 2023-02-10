package admin

import "context"

type Admin struct {
	adminUsername string
}

func (a *Admin) GetAdminCredentials() string {
	return a.adminUsername
}

func NewAdmin(ctx context.Context) (*Admin, error) {
	username, err := getAdminCredentials(ctx)
	if err != nil {
		return nil, err
	}

	return &Admin{
		adminUsername: username,
	}, nil

}

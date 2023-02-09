package domain

import "context"

type DB interface {
	GetClientAuthentication(ctx context.Context, cid, username, password string) (bool, error)
	GetClientAuthenticationACL(ctx context.Context, cid, topic, acc string) (bool, error)

	GetSuperuserAuthentication(ctx context.Context, username string) (bool, error)
}

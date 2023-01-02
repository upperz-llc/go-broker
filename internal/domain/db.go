package domain

import "context"

type DB interface {
	GetClientAuthentication(ctx context.Context, cid string) (bool, error)
}

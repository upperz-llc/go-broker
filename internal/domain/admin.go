package domain

import "context"

type Admin interface {
	IsAdmin(ctx context.Context, username string) bool
}

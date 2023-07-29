package hooks

import (
	"context"
	"errors"
	"os"

	rv8 "github.com/go-redis/redis/v8"

	"github.com/mochi-co/mqtt/v2/hooks/storage/redis"
)

func NewRedisPersistanceHookConfig(ctx context.Context) (*redis.Options, error) {
	rp, found := os.LookupEnv("REDIS_PASSWORD")
	if !found {
		return nil, errors.New("REDIS_PASSWORD")
	}

	return &redis.Options{
		Options: &rv8.Options{
			Addr:     "redis-10731.c253.us-central1-1.gce.cloud.redislabs.com:10731", // default redis address
			Password: rp,                                                             // your password
			DB:       0,                                                              // your redis db
			Username: "default",
		},
	}, nil
}

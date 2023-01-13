package firestore

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DB placeholder
type DB struct {
	DB     *firestore.Client
	Logger *log.Logger
}

// DeleteDevice placeholder
func (db *DB) GetClientAuthentication(ctx context.Context, cid string) (bool, error) {
	wr, err := db.DB.Collection("broker-auth").Doc(cid).Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return false, errors.New("something went wrong")
		}
	}

	if !wr.Exists() {
		return false, nil
	}

	enabled, err := wr.DataAt("enabled")
	if err != nil {
		return false, err
	}

	return enabled.(bool), nil
}

// GetClientAuthenticationACL placeholder
func (db *DB) GetClientAuthenticationACL(ctx context.Context, cid, topic string) (bool, error) {
	wr, err := db.DB.Collection("broker-auth").Doc(cid).Collection("acls").Doc(topic).Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return false, errors.New("something went wrong")
		}
	}

	if !wr.Exists() {
		return false, nil
	}

	allowed, err := wr.DataAt("allowed")
	if err != nil {
		db.Logger.Println(err)
		return false, err
	}
	return allowed.(bool), nil
}

func NewClient(ctx context.Context) (*DB, error) {
	return &DB{}, nil
}

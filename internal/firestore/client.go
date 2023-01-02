package firestore

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Firestore placeholder
type DB struct {
	DB *firestore.Client
}

// DeleteDevice placeholder
// func (db *DB) DeleteDevice(ctx context.Context, deviceID string) error {
// 	iter := db.DB.CollectionGroup("devices").Where("device_id", "==", deviceID).Documents(ctx)
// 	devices, err := iter.GetAll()
// 	if err != nil {
// 		return err
// 	}
// 	if len(devices) == 0 || len(devices) > 1 {
// 		return errors.New("test")
// 	}

// 	coliter := devices[0].Ref.Collections(ctx)

// 	for {
// 		doc, err := coliter.Next()
// 		if err == iterator.Done {
// 			break
// 		}

// 		fmt.Println(doc.Path)

// 		iter := doc.Documents(ctx)
// 		for {
// 			doc, err := iter.Next()
// 			if err == iterator.Done {
// 				break
// 			}
// 			if _, err := doc.Ref.Delete(ctx); err != nil {
// 				fmt.Println(err)
// 			}
// 		}
// 	}

// 	if _, err := devices[0].Ref.Delete(ctx); err != nil {
// 		return errors.New(fmt.Sprintf("Error : %s", err.Error()))
// 	}

// 	return nil
// }

// DeleteDevice placeholder
func (db *DB) GetClientAuthentication(ctx context.Context, cid string) (bool, error) {
	wr, err := db.DB.Collection("broker-auth").Doc(cid).Get(ctx)
	if err != nil {
		if status.Code(err) != codes.NotFound {
			return false, errors.New("something went wrong")
		}
	}

	return wr.Exists(), nil
}

// DeleteDevice placeholder
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
		fmt.Println(err)
		return false, err
	}
	return allowed.(bool), nil
}

func NewClient(ctx context.Context) (*DB, error) {
	return &DB{}, nil
}

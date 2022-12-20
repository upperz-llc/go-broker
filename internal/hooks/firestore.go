package hooks

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"

	firebase "firebase.google.com/go/v4"
)

type FirestoreAuthHook struct {
	mqtt.HookBase
	DB *firestore.Client
}

func (h *FirestoreAuthHook) ID() string {
	return "firebase-auth"
}

func (h *FirestoreAuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnACLCheck,
		mqtt.OnConnectAuthenticate,
	}, []byte{b})
}

func (h *FirestoreAuthHook) Init(config any) error {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	fdb, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	h.DB = fdb

	h.Log.Info().Msg("initialised")
	return nil
}

func (h *FirestoreAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	fmt.Println("OnConnectAuthenticate")

	return true
}

func (h *FirestoreAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	fmt.Println("OnACLCheck")

	return true
}

package hooks

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/patrickmn/go-cache"
	"github.com/upperz-llc/go-broker/internal/firestore"

	firebase "firebase.google.com/go/v4"
)

type FirestoreAuthHook struct {
	ACLCache *cache.Cache
	DB       *firestore.DB
	Logger   *log.Logger
	mqtt.HookBase
}

func (h *FirestoreAuthHook) ID() string {
	return "firebase-auth-hook"
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

	db, err := firestore.NewClient(context.Background())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	db.DB = fdb
	h.DB = db

	h.ACLCache = cache.New(5*time.Minute, 10*time.Minute)

	h.Logger.Println("initialized firestoreauthhook")
	return nil
}

func (h *FirestoreAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	h.Logger.Println("OnConnectAuthenticate")

	if found, err := h.DB.GetClientAuthentication(context.Background(), cl.ID); err != nil {
		return false
	} else {
		h.Logger.Printf("Connect Authenticate check result : %t\n", found)
		return found
	}

}

func (h *FirestoreAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	h.Logger.Println("OnACLCheck")

	allowed, found := h.ACLCache.Get(cl.ID + "-" + topic)
	if !found {
		h.Logger.Printf("Cache miss for ACL Check : %s for topic : %s\n", cl.ID, topic)
		if found, err := h.DB.GetClientAuthenticationACL(context.Background(), cl.ID, topic); err != nil {
			return false
		} else {
			h.Logger.Printf("ACL check result : %t\n", found)
			h.Logger.Printf("Adding ACL to cache for device id : %s for topic : %s\n", cl.ID, topic)
			if err := h.ACLCache.Add(cl.ID+"-"+topic, found, cache.DefaultExpiration); err != nil {
				h.Logger.Println(err)
			}

			return found
		}
	}

	h.Logger.Printf("Cache hit for ACL Check : %s for topic : %s\n", cl.ID, topic)
	return allowed.(bool)
}

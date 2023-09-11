package hooks

import (
	"bytes"
	"context"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/logging"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/patrickmn/go-cache"
	"github.com/upperz-llc/go-broker/internal/firestore"

	firebase "firebase.google.com/go/v4"
)

type FirestoreAuthHook struct {
	ACLCache                   *cache.Cache
	OnConnectAuthenticateCache *cache.Cache
	DB                         *firestore.DB
	Logger                     *logging.Logger
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
	h.OnConnectAuthenticateCache = cache.New(5*time.Minute, 10*time.Minute)

	h.Logger.StandardLogger(logging.Debug).Println("initialized firestoreauthhook")
	return nil
}

func (h *FirestoreAuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	h.Logger.StandardLogger(logging.Debug).Printf("OnConnectAuthenticate called for device : %s\n", cl.ID)

	// ADMIN CHECK
	if string(pk.Connect.Username) == "admin" && string(pk.Connect.Password) == "admin" {
		return true
	}

	allowed, found := h.OnConnectAuthenticateCache.Get(cl.ID)
	if !found {
		h.Logger.StandardLogger(logging.Debug).Printf("Cache miss for OnConnectAuthenticate for deviceID : %s\n", cl.ID)
		if allowed, err := h.DB.GetClientAuthentication(context.Background(), cl.ID, string(cl.Properties.Username)); err != nil {
			h.Logger.StandardLogger(logging.Error).Println(err)
			return false
		} else {
			h.Logger.StandardLogger(logging.Debug).Printf("Connect Authenticate check result : %t\n", allowed)
			if err := h.OnConnectAuthenticateCache.Add(cl.ID, allowed, cache.DefaultExpiration); err != nil {
				h.Logger.StandardLogger(logging.Debug).Printf("Failed to save OnConnectAuthenticate result to cache for device : %s\n", cl.ID)
				h.Logger.StandardLogger(logging.Error).Println(err)
				return allowed
			}
			return allowed
		}
	}

	h.Logger.StandardLogger(logging.Debug).Printf("Cache hit for OnConnectAuthenticate Check for device : %s\n", cl.ID)
	return allowed.(bool)
}

func (h *FirestoreAuthHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	h.Logger.StandardLogger(logging.Debug).Printf("OnACLCheck called for device : %s on topic : %s\n", cl.ID, topic)

	// ADMIN CHECK
	if string(cl.Properties.Username) == "admin" {
		return true
	}
	formatted_topic := strings.ReplaceAll(topic, "/", "")
	allowed, found := h.ACLCache.Get(cl.ID + "-" + formatted_topic)
	if !found {
		h.Logger.StandardLogger(logging.Debug).Printf("Cache miss for ACL Check : %s for topic : %s\n", cl.ID, topic)
		if found, err := h.DB.GetClientAuthenticationACL(context.Background(), cl.ID, formatted_topic); err != nil {
			return false
		} else {
			h.Logger.StandardLogger(logging.Debug).Printf("ACL check result : %t\n", found)
			h.Logger.StandardLogger(logging.Debug).Printf("Adding ACL to cache for device id : %s for topic : %s\n", cl.ID, formatted_topic)
			if err := h.ACLCache.Add(cl.ID+"-"+formatted_topic, found, cache.DefaultExpiration); err != nil {
				h.Logger.StandardLogger(logging.Error).Println(err)
			}

			return found
		}
	}

	h.Logger.StandardLogger(logging.Debug).Printf("Cache hit for ACL Check : %s for topic : %s\n", cl.ID, topic)
	return allowed.(bool)
}

package hooks

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/upperz-llc/go-broker/internal/domain"
	"github.com/upperz-llc/go-broker/internal/ps"
)

type GCPPubsubHook struct {
	Logger *log.Logger
	Pubsub ps.BrokerPubSub
	mqtt.HookBase
}

func (h *GCPPubsubHook) ID() string {
	return "gcp-pubsub-hook"
}

func (h *GCPPubsubHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
		mqtt.OnPublished,
	}, []byte{b})
}

func (h *GCPPubsubHook) Init(config any) error {
	ctx := context.Background()

	// Pull Env Variables
	pid, found := os.LookupEnv("GCP_PROJECT_ID")
	if !found {
		panic("GCP_PROJECT_ID not found ... panicing")
	}
	bt, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC")
	if !found {
		panic("BROKER_HOOK_GCPPUBSUB_TOPIC not found ... panicing")
	}

	// Create and configure logger
	lc, err := logging.NewClient(ctx, pid)
	if err != nil {
		panic("Failed to create client")
	}

	logger := lc.Logger("go-broker-log").StandardLogger(logging.Info)

	// Create and configure pubsub client
	pc, err := pubsub.NewClient(ctx, pid)
	if err != nil {
		panic(fmt.Errorf("pubsub.NewClient: %v", err))
	}

	topic := pc.Topic(bt)
	topic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	// Create internal broker logic
	bps := ps.BrokerPubSub{
		Logger: *logger,
		Topic:  topic,
	}

	h.Pubsub = bps
	h.Logger = logger

	h.Logger.Println("initialised gcp pubsub hook")
	return nil
}

func (h *GCPPubsubHook) OnConnect(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Msgf("client connected")
}

func (h *GCPPubsubHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnected")
}

func (h *GCPPubsubHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msgf("subscribed qos=%v", reasonCodes)
}

func (h *GCPPubsubHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msg("unsubscribed")
}

func (h *GCPPubsubHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Logger.Printf("Client %s published payload %s to client", cl.ID, string(pk.Payload))
	err := h.Pubsub.Publish("test", domain.MQTTEvent{
		Topic:   pk.TopicName,
		Payload: pk.Payload,
	})

	if err != nil {
		fmt.Println(err)
	}
}

// func Initialize(ctx context.Context) (*GCPPubsubHook, error) {

// }

package hooks

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/upperz-llc/go-broker/internal/domain"
	"github.com/upperz-llc/go-broker/internal/ps"
	internalpkg "github.com/upperz-llc/go-broker/pkg/domain"
)

type GCPPubsubHook struct {
	Logger         *logging.Logger
	Pubsub         domain.Pubsub
	publishTopic   *pubsub.Topic
	subscripeTopic *pubsub.Topic
	connectTopic   *pubsub.Topic

	mqtt.HookBase
}

func (h *GCPPubsubHook) ID() string {
	return "gcp-pubsub-hook"
}

func (h *GCPPubsubHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnPublished,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
	}, []byte{b})
}

func (h *GCPPubsubHook) Init(config any) error {
	ctx := context.Background()

	// Pull Env Variables
	pid, found := os.LookupEnv("GCP_PROJECT_ID")
	if !found {
		h.Logger.StandardLogger(logging.Debug).Println("GCP_PROJECT_ID not found")
	}
	bt, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC_PUBLISH")
	if !found {
		h.Logger.StandardLogger(logging.Debug).Println("BROKER_HOOK_GCPPUBSUB_TOPIC_PUBLISH not found")
	}
	st, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC_SUBSCRIBE")
	if !found {
		h.Logger.StandardLogger(logging.Debug).Println("BROKER_HOOK_GCPPUBSUB_TOPIC_SUBSCRIBE not found")
	}
	ct, found := os.LookupEnv("BROKER_HOOK_GCPPUBSUB_TOPIC_CONNECT")
	if !found {
		h.Logger.StandardLogger(logging.Debug).Println("BROKER_HOOK_GCPPUBSUB_TOPIC_CONNECT not found")
	}

	// Create and configure logger
	lc, err := logging.NewClient(ctx, pid)
	if err != nil {
		panic("Failed to create client")
	}

	logger := lc.Logger("go-broker-log")

	// Create and configure pubsub client
	pc, err := pubsub.NewClient(ctx, pid)
	if err != nil {
		panic(fmt.Errorf("pubsub.NewClient: %v", err))
	}

	pubslishtopic := pc.Topic(bt)
	pubslishtopic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	subscribetopic := pc.Topic(st)
	subscribetopic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	connecttopic := pc.Topic(ct)
	connecttopic.PublishSettings = pubsub.PublishSettings{
		DelayThreshold: 1 * time.Second,
		CountThreshold: 10,
	}

	// Create internal broker logic
	bps := ps.BrokerPubSub{
		Logger: *logger,
	}

	h.Pubsub = &bps
	h.publishTopic = pubslishtopic
	h.subscripeTopic = subscribetopic
	h.connectTopic = connecttopic
	h.Logger = logger

	h.Logger.StandardLogger(logging.Info).Println("initialised gcp pubsub hook")
	return nil
}

func (h *GCPPubsubHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Logger.StandardLogger(logging.Debug).Printf("Client %s unsubscribed to %s at %s", cl.ID, pk.TopicName, time.Now())
	if cl.ID == "admin" {
		return
	}
	err := h.Pubsub.Publish(h.subscripeTopic, internalpkg.MochiSubscribeMessage{
		ClientID:   cl.ID,
		Username:   string(cl.Properties.Username),
		Timestamp:  time.Now(),
		Subscribed: false,
		Topic:      pk.TopicName,
	})

	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}

}

func (h *GCPPubsubHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Logger.StandardLogger(logging.Debug).Printf("Client %s subscribed to %s with reason codes %s at %s", cl.ID, pk.TopicName, reasonCodes, time.Now())
	if cl.ID == "admin" {
		return
	}
	err := h.Pubsub.Publish(h.subscripeTopic, internalpkg.MochiSubscribeMessage{
		ClientID:   cl.ID,
		Username:   string(cl.Properties.Username),
		Timestamp:  time.Now(),
		Subscribed: true,
		Topic:      pk.TopicName,
	})

	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}

}

func (h *GCPPubsubHook) OnConnect(cl *mqtt.Client, pk packets.Packet) {
	h.Logger.StandardLogger(logging.Debug).Printf("Client %s connected at %s", cl.ID, time.Now())
	if cl.ID == "admin" {
		return
	}
	err := h.Pubsub.Publish(h.connectTopic, internalpkg.MochiConnectMessage{
		ClientID:  cl.ID,
		Username:  string(cl.Properties.Username),
		Timestamp: time.Now(),
		Connected: true,
	})

	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}
}

func (h *GCPPubsubHook) OnDisconnect(cl *mqtt.Client, connect_err error, expire bool) {
	h.Logger.StandardLogger(logging.Debug).Printf("Client %s disconnected at %s", cl.ID, time.Now())
	if cl.ID == "admin" {
		return
	}
	err := h.Pubsub.Publish(h.connectTopic, internalpkg.MochiConnectMessage{
		ClientID:  cl.ID,
		Username:  string(cl.Properties.Username),
		Timestamp: time.Now(),
		Connected: false,
	})

	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}
}

func (h *GCPPubsubHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Logger.StandardLogger(logging.Debug).Printf("Client %s published payload %s to client", cl.ID, string(pk.Payload))
	err := h.Pubsub.Publish(h.publishTopic, internalpkg.MochiPublishMessage{
		ClientID:  cl.ID,
		Topic:     pk.TopicName,
		Payload:   pk.Payload,
		Timestamp: time.Now(),
	})

	if err != nil {
		h.Logger.StandardLogger(logging.Error).Println(err)
	}
}

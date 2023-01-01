package hooks

import (
	"bytes"
	"fmt"
	"log"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/upperz-llc/go-broker/internal/domain"
	"github.com/upperz-llc/go-broker/internal/ps"
)

type ExampleHook struct {
	Logger *log.Logger
	Pubsub ps.BrokerPubSub
	mqtt.HookBase
}

func (h *ExampleHook) ID() string {
	return "events-example"
}

func (h *ExampleHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
		mqtt.OnPublished,
	}, []byte{b})
}

func (h *ExampleHook) Init(config any) error {
	h.Log.Info().Msg("initialised")
	return nil
}

func (h *ExampleHook) OnConnect(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Msgf("client connected")
}

func (h *ExampleHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnected")
}

func (h *ExampleHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msgf("subscribed qos=%v", reasonCodes)
}

func (h *ExampleHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msg("unsubscribed")
}

func (h *ExampleHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Logger.Printf("Client %s published payload %s to client", cl.ID, string(pk.Payload))
	err := h.Pubsub.Publish("test", domain.MQTTEvent{
		Topic:   pk.TopicName,
		Payload: pk.Payload,
	})

	if err != nil {
		fmt.Println(err)
	}
}

package ps

import (
	"context"
	"encoding/json"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/upperz-llc/go-broker/internal/domain"
)

// BrokerPubSub palceholder
type BrokerPubSub struct {
	Logger log.Logger
	Topic  *pubsub.Topic
}

func (bps *BrokerPubSub) Publish(me domain.MQTTEvent) error {
	ctx := context.Background()

	b, _ := json.Marshal(me)

	// TODO : store results to check later
	bps.Topic.Publish(ctx, &pubsub.Message{
		Data: b,
	})

	return nil
}

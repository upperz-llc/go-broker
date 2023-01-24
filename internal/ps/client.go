package ps

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
)

// BrokerPubSub palceholder
type BrokerPubSub struct {
	Logger logging.Logger
}

func (bps *BrokerPubSub) Publish(topic *pubsub.Topic, data interface{}) error {
	ctx := context.Background()
	b, _ := json.Marshal(data)

	// TODO : store results to check later
	topic.Publish(ctx, &pubsub.Message{
		Data: b,
	})

	return nil
}

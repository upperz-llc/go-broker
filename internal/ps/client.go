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
	b, _ := json.Marshal(data)

	// TODO : store results to check later
	topic.Publish(context.Background(), &pubsub.Message{
		Data: b,
	})

	return nil
}

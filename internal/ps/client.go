package ps

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/upperz-llc/go-broker/internal/domain"
)

// BrokerPubSub palceholder
type BrokerPubSub struct {
	Topic *pubsub.Topic
}

func (bps *BrokerPubSub) Publish(topic string, me domain.MQTTEvent) error {
	ctx := context.Background()
	var wg sync.WaitGroup
	var totalErrors uint64

	b, _ := json.Marshal(me)

	result := bps.Topic.Publish(ctx, &pubsub.Message{
		Data: b,
	})

	wg.Add(1)
	go func(res *pubsub.PublishResult) {
		defer wg.Done()
		// The Get method blocks until a server-generated ID or
		// an error is returned for the published message.
		id, err := res.Get(ctx)
		if err != nil {
			// Error handling code can be added here.
			fmt.Printf("Failed to publish: %v \n", err)
			atomic.AddUint64(&totalErrors, 1)
			return
		}
		fmt.Printf("Published message; msg ID: %v\n", id)
	}(result)

	wg.Wait()

	return nil
}

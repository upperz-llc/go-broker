package ps

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/upperz-llc/go-broker/internal/domain"
)

// BrokerPubSub palceholder
type BrokerPubSub struct {
	Logger log.Logger
	Topic  *pubsub.Topic
}

func (bps *BrokerPubSub) Publish(topic string, me domain.MQTTEvent) error {
	ctx := context.Background()
	var wg sync.WaitGroup
	var totalErrors uint64

	b, _ := json.Marshal(me)

	// TODO : store results to check later
	result := bps.Topic.Publish(ctx, &pubsub.Message{
		Data: b,
	})

	// now := time.Now()
	// mid, err := result.Get(ctx)
	// if err != nil {
	// 	bps.Logger.Printf("Failed to publish pubsub message with message ID %s with error %s\n", mid, err)
	// 	return err
	// }
	// bps.Logger.Printf("Published message with %s \n", mid)
	// bps.Logger.Println(time.Since(now))

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
		bps.Logger.Printf("Published message; msg ID: %v\n", id)
	}(result)

	wg.Wait()

	return nil
}

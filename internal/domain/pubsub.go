package domain

import "cloud.google.com/go/pubsub"

type Pubsub interface {
	Publish(topic *pubsub.Topic, data interface{}) error
}

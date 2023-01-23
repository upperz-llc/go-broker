package domain

import "time"

// MQTTEvent placeholder
type MQTTEvent struct {
	Topic   string `json:"topic"`
	Payload []byte `json:"payload"`
}

// MQTTEvent placeholder
type ConnectPayload struct {
	ClientID  string    `json:"client_id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Connected bool      `json:"connected"`
}

// MQTTEvent placeholder
type SubscribePayload struct {
	ClientID   string    `json:"client_id"`
	Username   string    `json:"username"`
	Topic      string    `json:"topic"`
	Subscribed bool      `json:"subscribed"`
	Timestamp  time.Time `json:"timestamp"`
}

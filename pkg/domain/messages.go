package domain

import "time"

// MochiPublishMessage placeholder
type MochiPublishMessage struct {
	ClientID  string    `json:"client_id"`
	Topic     string    `json:"topic"`
	Payload   []byte    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}

// MochiConnectMessage placeholder
type MochiConnectMessage struct {
	ClientID  string    `json:"client_id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Connected bool      `json:"connected"`
}

// MochiSubscribeMessage placeholder
type MochiSubscribeMessage struct {
	ClientID   string    `json:"client_id"`
	Username   string    `json:"username"`
	Topic      string    `json:"topic"`
	Subscribed bool      `json:"subscribed"`
	Timestamp  time.Time `json:"timestamp"`
}

package domain

// MQTTEvent placeholder
type MQTTEvent struct {
	Topic   string `json:"topic"`
	Payload []byte `json:"payload"`
}

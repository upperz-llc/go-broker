package api

type MessageRequest struct {
	Topic   string `json:"topic"`
	QoS     byte   `json:"qos"`
	Retain  bool   `json:"retain"`
	Payload any    `json:"payload,omitempty"`
}

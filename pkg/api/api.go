package api

import "encoding/json"

type MessageRequest struct {
	Topic   string          `json:"topic"`
	QoS     byte            `json:"qos"`
	Retain  bool            `json:"retain"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

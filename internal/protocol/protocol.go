package protocol

import go_json "github.com/goccy/go-json"

type MessageType int

const (
	MessageTypeSubmitName MessageType = 0
	MessageTypeEndGame    MessageType = 1
)

type Message struct {
	Type MessageType        `json:"type"`
	Data go_json.RawMessage `json:"data"`
}

type SubmitName struct {
	Name string `json:"name"`
}

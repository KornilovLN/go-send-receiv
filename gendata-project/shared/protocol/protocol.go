// protocol.go
package protocol

import (
	"encoding/json"
	"gendata-project/shared/types"
)

type Message struct {
	Type      string            `json:"type"`
	Timestamp int64             `json:"timestamp"`
	Header    types.BlockHeader `json:"header"`
	Data      []byte            `json:"data"`
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

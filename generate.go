package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type generateMessage struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	MessageID int    `json:"msg_id"`
	InReplyTo int    `json:"in_reply_to"`
}

func (s *server) generate(msg maelstrom.Message) error {
	// Unmarshal the message body as an loosely-typed map.
	var body generateMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Update the message type to return back.
	randNumber, err := generateUniqueID()
	if err != nil {
		return err
	}
	body.ID = randNumber
	body.Type = "generate_ok"

	// Echo the original message back with the updated message type.
	return s.node.Reply(msg, body)
}

func generateUniqueID() (int64, error) {
	var result int64
	err := binary.Read(rand.Reader, binary.BigEndian, &result)
	if err != nil {
		return 0, err
	}

	return result, err
}

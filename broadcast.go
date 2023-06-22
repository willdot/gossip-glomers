package main

import (
	"encoding/json"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type broadcast struct {
	ids    []int64
	idLock sync.Mutex

	topologies     map[string][]string
	topoligiesLock sync.Mutex
}

func newBroadcast() broadcast {
	return broadcast{
		ids:        make([]int64, 0),
		topologies: make(map[string][]string),
	}
}

type broadcastMessage struct {
	Message int64 `json:"message"`
}

func (s *server) broadcast(msg maelstrom.Message) error {
	var body broadcastMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.bc.idLock.Lock()
	defer s.bc.idLock.Unlock()

	s.bc.ids = append(s.bc.ids, body.Message)

	return s.node.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
}

func (s *server) read(msg maelstrom.Message) error {
	var ids []int64
	s.bc.idLock.Lock()
	ids = s.bc.ids
	s.bc.idLock.Unlock()

	return s.node.Reply(msg, map[string]any{
		"type":     "read_ok",
		"messages": ids,
	})
}

type topologyMessage struct {
	topologies map[string][]string `json:"topology"`
}

func (s *server) topology(msg maelstrom.Message) error {
	var body topologyMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.bc.topoligiesLock.Lock()
	s.bc.topologies = body.topologies
	s.bc.topoligiesLock.Unlock()

	return s.node.Reply(msg, map[string]any{
		"type": "topology_ok",
	})
}

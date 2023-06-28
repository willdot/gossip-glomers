package main

import (
	"encoding/json"
	"fmt"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type broadcast struct {
	idLock sync.Mutex
	idMap  map[int64]struct{}
}

func newBroadcast() broadcast {
	return broadcast{
		idMap: make(map[int64]struct{}),
	}
}

type broadcastMessage struct {
	Message   int64 `json:"message"`
	Propagate bool  `json:"propagate"`
}

func (s *server) broadcast(msg maelstrom.Message) error {
	var body broadcastMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.bc.idLock.Lock()
	s.bc.idMap[body.Message] = struct{}{}
	defer s.bc.idLock.Unlock()

	if body.Propagate {
		return nil
	}

	go s.propagateValueToOtherNodes(body.Message)

	return s.node.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})
}

func (s *server) propagateValueToOtherNodes(id int64) {
	for _, node := range s.node.NodeIDs() {
		if node == s.node.ID() {
			continue
		}

		err := s.node.Send(node, map[string]any{
			"type":      "broadcast",
			"propagate": true,
			"message":   id,
		})

		if err != nil {
			fmt.Printf("failed to propagate to node %s: %s", node, err)
		}
	}
}

func (s *server) read(msg maelstrom.Message) error {
	ids := make([]int64, 0, len(s.bc.idMap))

	s.bc.idLock.Lock()
	for id := range s.bc.idMap {
		ids = append(ids, id)
	}
	s.bc.idLock.Unlock()

	return s.node.Reply(msg, map[string]any{
		"type":     "read_ok",
		"messages": ids,
	})
}

func (s *server) topology(msg maelstrom.Message) error {
	return s.node.Reply(msg, map[string]any{
		"type": "topology_ok",
	})
}

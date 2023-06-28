package main

import (
	"encoding/json"
	"sync"
	"time"

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
	Message int64 `json:"message"`
}

func newServer(node *maelstrom.Node, bc broadcast) server {
	s := server{
		node: node,
		bc:   bc,
	}

	return s
}

func (s *server) broadcast(msg maelstrom.Message) error {
	var body broadcastMessage
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.node.Reply(msg, map[string]any{
		"type": "broadcast_ok",
	})

	if ok := s.checkIfIDExists(body.Message); ok {
		return nil
	}

	s.addID(body.Message)

	s.propagateValueToOtherNodes(msg.Src, body.Message)

	return nil
}

func (s *server) addID(id int64) {
	s.bc.idLock.Lock()
	defer s.bc.idLock.Unlock()

	s.bc.idMap[id] = struct{}{}
}

func (s *server) checkIfIDExists(id int64) bool {
	s.bc.idLock.Lock()
	defer s.bc.idLock.Unlock()

	_, ok := s.bc.idMap[id]

	return ok
}

func (s *server) propagateValueToOtherNodes(src string, id int64) {
	for _, n := range s.node.NodeIDs() {
		node := n
		if node == src || node == s.node.ID() {
			continue
		}

		ackd := make(chan struct{})
		for {
			select {
			case <-ackd:
				break
			default:
				body := map[string]any{
					"type":    "broadcast",
					"message": id,
				}

				s.node.RPC(node, body, func(msg maelstrom.Message) error {
					ackd <- struct{}{}
					return nil
				})

				time.Sleep(500 * time.Millisecond)
			}
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

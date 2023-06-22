package main

import (
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type server struct {
	node *maelstrom.Node

	bc broadcast
}

func main() {
	n := maelstrom.NewNode()

	server := server{
		node: n,
		bc:   newBroadcast(),
	}

	n.Handle("echo", server.echo)
	n.Handle("generate", server.generate)

	n.Handle("broadcast", server.broadcast)
	n.Handle("read", server.read)
	n.Handle("topology", server.topology)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

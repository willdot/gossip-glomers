maelstromCmd := ${HOME}/bin/maelstrom/maelstrom
gossipGlomers := ${GOPATH}/bin/gossip-glomers

build:
	@go install

1: build
	$(maelstromCmd) test -w echo --bin $(gossipGlomers) --node-count 1 --time-limit 10

2: build
	$(maelstromCmd) test -w unique-ids --bin $(gossipGlomers) --time-limit 30 --rate 1000 --node-count 3 --availability total --nemesis partition

3a: build
	$(maelstromCmd) test -w broadcast --bin $(gossipGlomers) --node-count 1 --time-limit 20 --rate 10

3b: build
	$(maelstromCmd) test -w broadcast --bin $(gossipGlomers) --node-count 5 --time-limit 20 --rate 10

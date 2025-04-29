build:
	go build

client: build
	./snek --mode=client

server: build
	./snek --mode=server

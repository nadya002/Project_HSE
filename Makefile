GO := $(shell find . -name "*.go")

all: server client

server: bin/server

bin/server: $(GO)
	go build -o bin/server cmd/server.go

client: bin/client

bin/client: $(GO)
	go build -o bin/client cmd/client.go

.PHONY: clean docker.build docker.build.client docker.build.server

docker.build: docker.build.client docker.build.server

docker.build.client:
	docker build -f docker/client/Dockerfile . -t audio-client:latest

docker.build.server:
	docker build -f docker/server/Dockerfile . -t audio-server:latest

clean:
	rm -rf bin
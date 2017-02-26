GOBUILD := go build -v
export GOARCH = amd64

.PHONY: all clean client server

all: client server

ifeq ($(OS), linux)
export GOOS=linux
else
export GOOS=darwin
endif

client: build/client

server: build/server

build/client: client/client.go
	$(GOBUILD) -o $@ $<

build/server: server/server.go
	$(GOBUILD) -o $@ $<

clean:
	rm -rf build

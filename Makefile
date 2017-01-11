GOBUILD := go build -v

.PHONY: all clean client server

all: build build/client_x64 build/server_x64

client: build/client_x64

server: build/server_x64

build/client_x64: client/client.go
	$(GOBUILD) -o $@ $<

build/server_x64: server/server.go
	$(GOBUILD) -o $@ $<

clean:
	rm build/client* build/server*

build:
	@mkdir -p $@

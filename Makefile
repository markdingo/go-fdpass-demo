all: client server

client: client.go common.go
	go build -o client client.go common.go

server: server.go common.go
	go build -o server server.go common.go

clean:
	rm -f client server

.PHONY: fmt
fmt:
	find . -name '*.go' -type f -print | xargs gofmt -s -w

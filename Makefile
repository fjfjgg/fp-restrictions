export GO111MODULE=on

all: fp-restrictions

fp-restrictions:
	go build -o ./bin/fp-restrictions ./cmd/fp-restrictions

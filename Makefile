BINARIES := api indexer
# MAKEFLAGS += --no-print-directory

.PHONY: clean build $(BINARIES)

build:
	go build ./...

clean:
	go clean ./...
	find . -type f -name 'swagger.json' -exec rm -f {} +

test:
	go test ./...

run-indexer: build
	go run cmd/indexer/main.go

run-api: build
	go run cmd/api/main.go
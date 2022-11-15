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

swagger:
	swagger generate spec -o ./swagger.yaml --scan-models
swagger-serve: swagger
	swagger serve -F=swagger swagger.yaml

run-indexer: build
	go run cmd/indexer/main.go --env=./docker/dev/local.env

run-api: build
	go run cmd/api/main.go --env=./docker/dev/local.env
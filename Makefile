BINARIES := api indexer
IMAGE=directory
TAG=latest

# MAKEFLAGS += --no-print-directory

.PHONY: clean build $(BINARIES)

build:
	go build ./...

clean:
	go clean ./...
	find . -type f -name 'swagger.json' -exec rm -f {} +

test:
	go test ./...

test-unit:
	go test -v -short ./...

swagger:
	swagger generate spec -o ./docs/swagger.yaml --scan-models

# redoc-cli: install with `npm install -g redoc-cli`
swagger-html: swagger
	redoc-cli bundle -o docs/swagger.html docs/swagger.yaml

swagger-serve: swagger
	swagger serve -F=swagger swagger.yaml

run-indexer: build
	go run cmd/indexer/main.go --env=./docker/dev/local.env

run-api: build
	go run cmd/api/main.go --env=./docker/dev/local.env

db-migrate:
	tern migrate -c db/tern.conf -m db

lint:
	@./scripts/lint.sh

install:
	go install ./cmd/...

docker-build: swagger-html
	@docker build --platform=linux/amd64 . --file Dockerfile -t ${IMAGE}:${TAG}

docker-tag:
	@docker tag ${IMAGE}:${TAG} ghcr.io/arkeonetwork/${IMAGE}:${TAG}

docker-push:
	@docker push ghcr.io/arkeonetwork/${IMAGE}:${TAG}

push-image: docker-build docker-tag docker-push

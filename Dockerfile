#
# Arkeo Directory
#

ARG GO_VERSION="1.19"

#
# Build
#
FROM golang:${GO_VERSION} as builder

ARG GIT_VERSION
ARG GIT_COMMIT

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TAG=latest
RUN make install

# FROM node:19-bullseye as docs
# docs
FROM builder as docs
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y curl apt-transport-https gnupg debian-keyring debian-archive-keyring nodejs npm

RUN curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/gpg.2F8CB673971B5C9E.key' | gpg --dearmor -o /usr/share/keyrings/go-swagger-go-swagger-archive-keyring.gpg
RUN curl -1sLf 'https://dl.cloudsmith.io/public/go-swagger/go-swagger/config.deb.txt?distro=debian&codename=any-version' > /etc/apt/sources.list.d/go-swagger-go-swagger.list
RUN apt update -y
# https://github.com/go-swagger/go-swagger/releases/download/v0.30.4/swagger_linux_amd64 linux/amd64 binary if above fails again
RUN apt-get install -y swagger
RUN npm install -g redoc-cli

COPY --from=builder /app /app

WORKDIR /app
RUN make swagger-html

# final
FROM ubuntu:kinetic

RUN apt-get update -y && \
    apt-get upgrade -y && \
    apt-get install -y jq curl htop vim ca-certificates
RUN update-ca-certificates

# clean up
RUN apt-get clean && \
      rm -rf /var/lib/apt/lists/*

# binaries
COPY --from=builder /go/bin/indexer /go/bin/api /usr/bin/
COPY --from=docs /app/docs/swagger.html /var/www/html/index.html
COPY --from=docs /app/docs/swagger.yaml /var/www/html/swagger.yaml

COPY scripts /scripts

WORKDIR /root
RUN rm -rf /app
CMD ["indexer"]

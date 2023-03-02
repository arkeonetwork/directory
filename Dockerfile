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

# docs
FROM golang:${GO_VERSION} as docs
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y curl nodejs npm

RUN curl -sLf https://github.com/go-swagger/go-swagger/releases/download/v0.30.4/swagger_linux_amd64 -o /usr/local/bin/swagger
RUN chmod +x /usr/local/bin/swagger

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

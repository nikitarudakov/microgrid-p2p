FROM golang:1.24-alpine AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd/inventory/ .
COPY ./internal ./internal

RUN go build -o inventory-service

CMD ["/build/inventory-service"]
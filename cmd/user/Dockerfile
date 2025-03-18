FROM golang:1.24-alpine AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd/user/ .
COPY ./internal ./internal
COPY ./pkg/services/user ./pkg/services/user

RUN go build -o user-service

CMD ["/build/user-service"]
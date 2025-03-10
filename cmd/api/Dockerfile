FROM golang:1.24-alpine AS base

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd/api/ .
COPY ./api ./api
COPY ./internal ./internal

RUN go build -o graphql-api

CMD ["/build/graphql-api"]
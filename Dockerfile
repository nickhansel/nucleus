ARG GO_VERSION=1.19

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
# copy env
RUN go build -o ./app ./cmd/main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/app .
COPY --from=builder /api/.env .
# move the .env file to ../.env of the app
RUN mv .env ../.env

EXPOSE 8080

ENTRYPOINT ["./app"]
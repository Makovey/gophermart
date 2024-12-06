FROM golang:1.22.5-alpine AS builder
USER root

RUN apk add --upgrade --no-cache ca-certificates && update-ca-certificates

COPY . /github.com/Makovey/gophermart
WORKDIR /github.com/Makovey/gophermart

RUN go mod download
RUN go build -o ./bin/gophermart cmd/gophermart/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder ./bin/gophermart .

CMD ["./gophermart"]
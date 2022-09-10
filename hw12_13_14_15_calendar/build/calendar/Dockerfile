# syntax=docker/dockerfile:1
FROM golang:1.19-alpine AS builder
LABEL stage=builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o calendar cmd/calendar/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calendar /app/calendar

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

EXPOSE 8080 50051

CMD ./calendar -config ${CONFIG_FILE}
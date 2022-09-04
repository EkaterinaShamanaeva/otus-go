# syntax=docker/dockerfile:1
FROM golang:1.19-alpine AS builder
LABEL stage=builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o calendar cmd/calendar/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/calendar /app/calendar

EXPOSE 8888 50051

ENTRYPOINT ["./calendar"]
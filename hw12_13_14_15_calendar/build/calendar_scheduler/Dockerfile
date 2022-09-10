#FROM golang:1.19-alpine AS builder
#LABEL stage=builder
#
#WORKDIR /app
#
#COPY go.mod ./
#COPY go.sum ./
#
#RUN go mod download
#
#COPY . .
#
#RUN go build -o calendar_scheduler cmd/calendar_scheduler/main.go
#
#FROM alpine:latest
#
#WORKDIR /app
#
#COPY --from=builder /app/calendar_scheduler /app/calendar_scheduler
#
#ENTRYPOINT ["./calendar_scheduler"]
# syntax=docker/dockerfile:1
# Собираем в гошке
FROM golang:1.19 as build

ENV BIN_FILE /opt/calendar/calendar-scheduler
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar_scheduler/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="calendar_scheduler"
LABEL MAINTAINERS="student@otus.ru"

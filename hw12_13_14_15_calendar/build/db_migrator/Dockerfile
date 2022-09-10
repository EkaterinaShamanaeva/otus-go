#syntax=docker/dockerfile:1
FROM golang:1.19 as build
LABEL stage=builder

RUN go install github.com/pressly/goose/cmd/goose@latest

WORKDIR /app

RUN go build -o cmd/db_migrator/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/db_migrator /app/db_migrator

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ./db_migrator -config ${CONFIG_FILE}
#FROM golang:1.19 as build
#
#ENV BIN_FILE /opt/calendar/db_migrator
#ENV CODE_DIR /go/src/
#
#ADD . /go/src/db_migrator
#WORKDIR ${CODE_DIR}
#
## Кэшируем слои с модулями
#RUN go install github.com/pressly/goose/cmd/goose@latest
#COPY . .
##COPY . ${CODE_DIR}
#
## Собираем статический бинарник Go (без зависимостей на Си API),
## иначе он не будет работать в alpine образе.
#ARG LDFLAGS
#RUN CGO_ENABLED=0 go build \
#        -ldflags "$LDFLAGS" \
#        -o ${BIN_FILE} cmd/db_migrator/main.go
#
## На выходе тонкий образ
#FROM alpine:3.9
#
#LABEL ORGANIZATION="OTUS Online Education"
#LABEL SERVICE="db_migrator"
#LABEL MAINTAINERS="student@otus.ru"
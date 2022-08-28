BINDIR := "./bin"
BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
MIGRATOR := "$(BINDIR)/db_migrator"
SCHEDULER := "$(BINDIR)/calendar_scheduler"
SENDER := "$(BINDIR)/calendar_sender"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(MIGRATOR) -ldflags "$(LDFLAGS)" ./cmd/db_migrator
	go build -v -o $(SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

run: build
	$(BIN) -config ./configs/config.yaml
	$(SCHEDULER) -config ./configs/config_schedule.yaml
	$(SENDER) -config ./configs/config_rmq.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

migrate:
	$(MIGRATOR) -config ./configs/config.yaml -dir=./migrations up

generate:
	go generate ./...

.PHONY: build run build-img run-img version test lint

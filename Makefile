# go
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOINSTALL=$(GOCMD) install
MAKE=make

# binary
BUILD_PATH=build
SERVICE=./cmd/service
SERVICE_BINARY=server

# Docker
DOCKER=docker
DC=docker compose
DCUP=up -d

default: run

.PHONY: test
test:
	$(GOTEST) -v ./...

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BUILD_PATH)

.PHONY: build
build: 
	$(GOBUILD) -o $(BUILD_PATH)/$(SERVICE_BINARY) -v $(SERVICE)

.PHONY: install
install:
	$(GOINSTALL) $(SERVICE)

.PHONY: docs
docs:
	# more recent version not working (1.7.0). had to downgrade to 1.6.7 in order to make it work
	# swag init --parseDependency --parseInternal -g cmd/service/main.go
	swag init -g cmd/service/main.go

.PHONY: run
run: docs
	reflex -r '\.go$$' -R './docs/*.go' -s -- sh -c 'config="./config/local.json" $(GORUN) $(SERVICE)'

up:
	$(DC) $(DCUP) service

up-asynq:
	$(DC) $(DCUP) --remove-orphans redis asynqmon

up-kafka:
# MY_IP=192.168.1.12 $(DC) $(DCUP) --remove-orphans zk1 zk2 zk3 kafka1 kafka2 kafka3
	$(DC) $(DCUP) --remove-orphans zookeeper kafka

up-nats:
	$(DC) $(DCUP) --remove-orphans nats1 nats2 nats3
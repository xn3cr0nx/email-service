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
DC=docker-compose
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

docker-server:
	$(DC) $(DCUP) service
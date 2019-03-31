TARGET := metricplower
TARGET_PATH := cmd/metricplower/main.go
BIN_PATH := bin

VERSION := $(shell sh -c 'git describe --always --tags')
DOCKER_VERSION := vlamug/$(TARGET):$(VERSION)

GO := go
DOCKER := docker
GOOS ?= linux
GOARCH ?= amd64

build:
	@echo ">>> building binary..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(TARGET) $(TARGET_PATH)
	mv $(TARGET) $(BIN_PATH)

docker-build: build
	@echo ">>> building docker image..."
	$(DOCKER) build -t $(TARGET):$(VERSION) .

docker-push: docker-build
	@echo ">>> logging on docker hub..."
	$(DOCKER) login
	@echo ">>> tagging image as '$(DOCKER_VERSION)'..."
	$(DOCKER) tag $(TARGET):$(VERSION) $(DOCKER_VERSION)
	@echo ">>> pushing image..."
	$(DOCKER) push $(DOCKER_VERSION)

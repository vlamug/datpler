TARGET := ratibor
TARGET_PATH := cmd/ratibor/main.go
BIN_PATH := bin

VERSION := $(shell sh -c 'git describe --always --tags')
DOCKER_VERSION := vlamug/$(TARGET):$(VERSION)

GO := go
DOCKER := docker
GOOS ?= linux
GOARCH ?= amd64

LOG_PATH ?= /var/log/ratibor/ratibor.logs

build:
	@echo ">>> building binary ..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build -o $(TARGET) $(TARGET_PATH)
	mv $(TARGET) $(BIN_PATH)

docker-build: build
	@echo ">>> building docker image ..."
	$(DOCKER) build -t $(TARGET):$(VERSION) .

docker-push: docker-build
	@echo ">>> logging on docker hub ..."
	$(DOCKER) login
	@echo ">>> tagging image as '$(DOCKER_VERSION)'..."
	$(DOCKER) tag $(TARGET):$(VERSION) $(DOCKER_VERSION)
	@echo ">>> pushing image..."
	$(DOCKER) push $(DOCKER_VERSION)

run: build
	@echo ">>> running app with logs in $(LOG_PATH) ..."
	./$(BIN_PATH)/$(TARGET) > $(LOG_PATH) 2>&1 &

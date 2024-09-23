DOCKER := $(shell if which podman >/dev/null 2>/dev/null; then echo podman; else echo docker; fi)

.DEFAULT_GOAL := build

BUILD := build

HOST_ARCH := $(shell arch)
ifeq ($(HOST_ARCH),x86_64)
	HOST_ARCH := amd64
endif

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
	HOST_OS := linux
endif
ifeq ($(UNAME_S),Darwin)
	HOST_OS := darwin
endif

VERSION := 0.0.1

GO_SOURCES := \
	../go.mod \
	../go.sum \
	$(shell find . -name \*.go)

$(BUILD)/htc: $(GO_SOURCES)
	@mkdir -p $(dir $@)
	go build -o $@ ./cmd/htc

.PHONY: build
build: $(BUILD)/htc

.PHONY: test
test:
	go test ./...

DOCKER := $(shell if which podman >/dev/null 2>/dev/null; then echo podman; else echo docker; fi)
BUILD := build
DIST_ARCH := $(shell if [ "$$(uname -m)" = "x86_64" ]; then echo amd64; else echo arm64; fi)
VERSION := 0.0.1
DIST_TGZ := $(BUILD)/htccli-$(VERSION)-$(DIST_ARCH).tar.gz

BUILD_OPTS := --platform linux/$(DIST_ARCH)

GO_SOURCES := \
	go.mod \
	go.sum \
	$(shell find . -name \*.go)

GO_LINUX_BINARIES := \
	$(BUILD)/htccli.linux-amd64 \
	$(BUILD)/htccli.linux-arm64

GO_LINUX_ARCHIVES := \
	$(BUILD)/htccli.linux-amd64.tar.gz \
    $(BUILD)/htccli.linux-arm64.tar.gz

IMAGE_NAME := htc_storage_cli

.PHONY: format
format:
	test -z $(gofmt -l .)

# Pattern for building all architectures.
# E.g. build/htccli.linux-amd64
$(BUILD)/htccli.%: $(GO_SOURCES)
	@mkdir -p $(BUILD)
	CGO_ENABLED=0 \
		GOOS=$(firstword $(subst -, ,$*)) \
		GOARCH=$(lastword $(subst -, ,$*)) \
		go build -o $@

.PHONY: image
image: $(GO_LINUX_BINARIES)
	@$(DOCKER) manifest rm $(IMAGE_NAME):$(VERSION) || true
	$(DOCKER) manifest create $(IMAGE_NAME):$(VERSION)
	$(DOCKER) build \
		--manifest $(IMAGE_NAME):$(VERSION) \
		--platform linux/amd64 \
		--build-arg BINARY_FILE=$(BUILD)/htccli.linux-amd64 .
	$(DOCKER) build \
		--manifest $(IMAGE_NAME):$(VERSION) \
		--platform linux/arm64 \
		--build-arg BINARY_FILE=$(BUILD)/htccli.linux-arm64 .


#$(BUILD)/htccli.%.tar.gz: $(BUILD)/htccli.%
#	@mkdir -p $(BUILD)/dist.$*
#	cp $< $(BUILD)/dist.$*/htccli
#	tar -czf $@ -C $(BUILD)/dist.$* .

.PHONY: archive
archive: $(GO_LINUX_ARCHIVES)


.PHONY: build-binary
build-binary: format
	go build -o ./htccli

.PHONY: test
test:
	go test

.PHONY: clean
clean:
	rm -f htccli

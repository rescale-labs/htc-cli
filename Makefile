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

$(BUILD)/htccli.linux-amd64: $(GO_SOURCES)
	@mkdir -p $(BUILD)
	CGO_ENABLED=0 \
		GOOS=linux \
		GOARCH=amd64 \
		go build -o $@

$(BUILD)/htccli.linux-arm64: $(GO_SOURCES)
	@mkdir -p $(BUILD)
	CGO_ENABLED=0 \
		GOOS=linux \
		GOARCH=arm64 \
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


$(BUILD)/htccli.linux-arm64.tar.gz: $(BUILD)/htccli.linux-arm64
	@mkdir -p $(BUILD)/dist.$*
	cp $< $(BUILD)/dist.$*/htccli
	tar -czf $@ -C $(BUILD)/dist.$* .

$(BUILD)/htccli.linux-amd64.tar.gz: $(BUILD)/htccli.linux-amd64
	@mkdir -p $(BUILD)/dist.$*
	cp $< $(BUILD)/dist.$*/htccli
	tar -czf $@ -C $(BUILD)/dist.$* .

.PHONY: archive
archive: $(GO_LINUX_ARCHIVES)

.PHONY: binary
binary: format
	go build -o ./htccli

.PHONY: test
test:
	go test && cd cli && go test

.PHONY: clean
clean:
	rm -rf $(BUILD)


.PHONY: push-dev
push-dev: $(GO_LINUX_BINARIES)
	$(DOCKER) build \
			--file Dockerfile.alpine \
    		-t us-docker.pkg.dev/rescale-htc-scale-test-sandbox/shared/$(IMAGE_NAME)_amd64:$(VERSION) \
    		--platform linux/amd64 \
    		--build-arg BINARY_FILE=$(BUILD)/htccli.linux-amd64 .
	$(DOCKER) build \
			--file Dockerfile.alpine \
    		-t us-docker.pkg.dev/rescale-htc-scale-test-sandbox/shared/$(IMAGE_NAME)_arm64:$(VERSION) \
    		--platform linux/arm64 \
    		--build-arg BINARY_FILE=$(BUILD)/htccli.linux-arm64 .
	$(DOCKER) push us-docker.pkg.dev/rescale-htc-scale-test-sandbox/shared/$(IMAGE_NAME)_amd64:$(VERSION)
	$(DOCKER) push us-docker.pkg.dev/rescale-htc-scale-test-sandbox/shared/$(IMAGE_NAME)_arm64:$(VERSION)

.PHONY: push-prod
push-prod: $(GO_LINUX_BINARIES)
	$(DOCKER) build \
			--file Dockerfile.alpine \
    		-t us-docker.pkg.dev/rescale-htc-prod/shared/$(IMAGE_NAME)_amd64:$(VERSION) \
    		--platform linux/amd64 \
    		--build-arg BINARY_FILE=$(BUILD)/htccli.linux-amd64 .
	$(DOCKER) build \
			--file Dockerfile.alpine \
    		-t us-docker.pkg.dev/rescale-htc-prod/shared/$(IMAGE_NAME)_arm64:$(VERSION) \
    		--platform linux/arm64 \
    		--build-arg BINARY_FILE=$(BUILD)/htccli.linux-arm64 .
	$(DOCKER) push us-docker.pkg.dev/rescale-htc-prod/shared/$(IMAGE_NAME)_amd64:$(VERSION)
	$(DOCKER) push us-docker.pkg.dev/rescale-htc-prod/shared/$(IMAGE_NAME)_arm64:$(VERSION)

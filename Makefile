DOCKER := $(shell if which podman >/dev/null 2>/dev/null; then echo podman; else echo docker; fi)
BUILD := build
DIST_ARCH := $(shell if [ "$$(uname -m)" = "x86_64" ]; then echo amd64; else echo arm64; fi)
BUILD_OPTS := --platform linux/$(DIST_ARCH)
VERSION := $(shell cat VERSION)
DIST_TGZ := $(BUILD)/htccli-$(VERSION)-$(DIST_ARCH).tar.gz

.PHONY: format
format:
	test -z $(gofmt -l .)

.PHONY: image
image: format
	$(DOCKER) build -t htc_storage_cli_$(DIST_ARCH):$(VERSION) $(BUILD_OPTS) .

$(DIST_TGZ): image
	@mkdir -p $(BUILD)
	$(DOCKER) rm htccli-archive || true
	$(DOCKER) create --platform linux/$(DIST_ARCH) \
    		--name htccli-archive \
    		--entrypoint /opt/rescale/bin/htccli \
    		htc_storage_cli_$(DIST_ARCH):$(VERSION)
	$(DOCKER) export htccli-archive | gzip -c > $@.tmp
	$(DOCKER) rm htccli-archive
	mv $@.tmp $@

.PHONY: archive
archive: $(DIST_TGZ)

.PHONY: build-binary
build-binary: format
	go build -o ./htccli

.PHONY: test
test:
	go test

.PHONY: clean
clean:
	rm -f htccli

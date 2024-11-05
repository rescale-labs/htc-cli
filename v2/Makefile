.DEFAULT_GOAL := build

# Manually update this when tagging a new release.
VERSION := v0.0.1

.PHONY: echo-version
echo-version:
	@echo $(VERSION)

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

# Create statically linked binaries.
GO_OPTS := CGO_ENABLED=0

GO_SOURCES := \
	../go.mod \
	../go.sum \
	$(shell find . -name \*.go)

#
# api/_oas generated code
#

# Update original swagger.json from HTC prod API.
.PHONY: refresh-swagger
refresh-swagger:
	curl -o api/swagger.json -H "Accept: application/json" \
    	https://htc.rescale.com/q/openapi

# Patches swagger.json so that it has the types and response
# codes, etc. that we need in our generated code.
api/swagger-patched.json: api/swagger.json api/swagger.jsonnet
	cd $(dir $@) && jsonnet swagger.jsonnet -o $(notdir $@)

$(BUILD)/swagger-sorted.json: api/swagger.json
	jsonnet -e "import '$<'" > $@

# Convenience method for diffing swagger.
.PHONY: diff-swagger
diff-swagger: $(BUILD)/swagger-sorted.json api/swagger-patched.json
	@echo "\nDifferences from original and swagger-patched.json:"
	diff -u $^ || true # diff exits 1 for differences

# Updates ogen OpenAPI generated code.
.PHONY: oapi
oapi: api/swagger-patched.json
	go generate ./api/...

#
# Static binaries for every target architecture.
#

# Pattern rule for any individual arch.
#
# Follows pattern of build/dist/$GOOS/$GOARCH/htc
$(BUILD)/dist/%/rescale-htc/htc: $(GO_SOURCES)
	@mkdir -p $(dir $@)
	GOOS=$(subst /,,$(dir $*)) GOARCH=$(notdir $*) $(GO_OPTS) go build -o $@ ./cmd/htc

# Overridding static rule for OS X universal. Uses a portable version of
# lipo.
$(BUILD)/dist/darwin/universal/rescale-htc/htc: $(foreach arch,arm64 amd64,$(BUILD)/dist/darwin/$(arch)/rescale-htc/htc)
	@mkdir -p $(dir $@)
	go run github.com/randall77/makefat@7ddd0e42c8442593c87c1705a5545099604008e5 \
		$@ $<

# Convenience target for building+copying host OS binary to build/htc.
$(BUILD)/htc: $(BUILD)/dist/$(HOST_OS)/$(HOST_ARCH)/rescale-htc/htc
	@mkdir -p $(dir $@)
	@cp $< $@

.PHONY: build
build: $(BUILD)/htc

DIST_BINARIES := $(foreach \
	arch,arm64 amd64,$(foreach \
		os,darwin linux,$(BUILD)/dist/$(os)/$(arch)/rescale-htc/htc)) \
	$(BUILD)/dist/darwin/universal/htc

ARCHIVE_PREFIX := $(BUILD)/dist/rescale-htc-$(VERSION)

DIST_ARCHIVES := $(foreach \
	arch,arm64 amd64,$(ARCHIVE_PREFIX)-linux-$(arch).tar.gz) \
	$(ARCHIVE_PREFIX)-darwin-universal.tar.gz

define tar-archive
cp CHANGES.md $(dir $<)/
tar -czf $@ -C $(dir $<).. rescale-htc
endef

$(ARCHIVE_PREFIX)-linux-amd64.tar.gz: $(BUILD)/dist/linux/amd64/rescale-htc/htc
	$(tar-archive)

$(ARCHIVE_PREFIX)-linux-arm64.tar.gz: $(BUILD)/dist/linux/arm64/rescale-htc/htc
	$(tar-archive)

$(ARCHIVE_PREFIX)-darwin-universal.tar.gz: $(BUILD)/dist/darwin/universal/rescale-htc/htc
	$(tar-archive)

.PHONY: dist
dist: $(DIST_ARCHIVES)

#
# Tests and docs
#

.PHONY: test
test:
	$(GO_OPTS) go test ./...

godoc:
	@echo "view oapi docs on: http://localhost:6060/pkg/github.com/rescale/htc-storage-cli/v2/api/_oas/"
	@echo
	godoc -http=localhost:6060

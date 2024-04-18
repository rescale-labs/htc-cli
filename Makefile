DOCKER := $(shell if which podman >/dev/null 2>/dev/null; then echo podman; else echo docker; fi)
BUILD_OPTS := --platform linux/amd64

build-container:
	$(DOCKER) build -t htc-storage-cli $(BUILD_OPTS) .

build-binary:
	go build -o ./htccli

clean:
	rm -f htccli

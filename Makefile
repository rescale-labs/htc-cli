DOCKER := $(shell if which podman >/dev/null 2>/dev/null; then echo podman; else echo docker; fi)
BUILD_OPTS := --platform linux/amd64

format:
	test -z $(gofmt -l .)

build-container: format
	$(DOCKER) build -t htc-storage-cli $(BUILD_OPTS) .

build-binary: format
	go build -o ./htccli

test:
	go test

clean:
	rm -f htccli

IMAGE              := databus23/concourse-swift-resource

ifneq ($(http_proxy),)
BUILD_ARGS+= --build-arg http_proxy=$(http_proxy) --build-arg https_proxy=$(https_proxy) --build-arg no_proxy=$(no_proxy)
endif

build: export GOOS=linux
build: export CGO_ENABLED=0
build:
	go build -o bin/check ./cmd/check
	go build -o bin/in ./cmd/in
	go build -o bin/out ./cmd/out

.PHONY: test
test:
	go vet ./cmd/... ./pkg/...
	go test -v ./cmd/... ./pkg/...

image:
	docker build -t $(IMAGE) $(BUILD_ARGS) .

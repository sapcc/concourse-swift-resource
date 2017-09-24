IMAGE := meshcloud/concourse-swift-resource
TAG   := 1.3.1

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

.PHONY: integration
integration:
	docker build -t $(IMAGE)-test -f Dockerfile.tdd .
	docker run $(IMAGE)-test

image:
	docker build -t $(IMAGE):$(TAG) $(BUILD_ARGS) .

release: image
	docker tag $(IMAGE):$(TAG) $(IMAGE):latest
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest

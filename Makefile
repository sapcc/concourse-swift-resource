BUILD_IMAGE:=docker.mo.sap.corp/monsoon/arc-build
IMAGE=docker.mo.sap.corp/concourse/swift-resource
build:
	docker run --rm -v $(CURDIR):/build -w /build $(BUILD_IMAGE) gb build -f -ldflags="-w -s"
	docker build --rm -t $(IMAGE) .

.PHONY: test
test:
	go vet ./src/...
	golint ./src/...
	gb test

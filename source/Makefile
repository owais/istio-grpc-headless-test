IMAGE ?= owais/istio-grpc-test
VERSION ?= dev-$$(git rev-parse --short HEAD)

.PHONY: docker-setup
docker-setup:
	docker buildx create --use

.PHONY: docker-build
docker-build:
		docker buildx build -t ${IMAGE}:${VERSION} --target final --platform=linux/arm64,linux/amd64 .

.PHONY: docker-push
docker-push:
		docker buildx build --push -t ${IMAGE}:${VERSION} --target final --platform=linux/arm64,linux/amd64 .

.PHONY: docker-build-local
docker-build-local:
		DOCKER_BUILDKIT=1 docker build -t istio-grpc-test --target final .
		docker tag istio-grpc-test ${IMAGE}:${VERSION}

.PHONY: docker-push
docker-push-local-image:
	docker push ${IMAGE}:${VERSION}

.PHONY: run-server
run-server:
	go run . server

.PHONY: run-client
run-client:
	go run . client 

.PHONY: build 
build:
	go build -o ./bin/istio-grpc-test .
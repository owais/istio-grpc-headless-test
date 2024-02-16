IMAGE ?= owais/istio-grpc-test
VERSION ?= dev-$$(git rev-parse --short HEAD)

REALM?=
NAMESPACE?=istio-grpc-test
SERVICE_NAME?=istio-grpc-test
KUBE_AUTH_METHOD?=local


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
		DOCKER_BUILDKIT=1 docker build -t ${SERVICE_NAME} --target final .
		docker tag ${SERVICE_NAME} ${IMAGE}:${VERSION}

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

check-env:
ifeq ($(SERVICE_NAME),)
	$(error SERVICE_NAME is undefined)
endif
ifeq ($(REALM),)
	$(error REALM is undefined)
endif
ifeq ($(NAMESPACE),)
	$(error NAMESPACE is undefined)
endif


deploy: check-env
	export NAMESPACE=$(NAMESPACE) REALM=$(REALM) SERVICE_NAME=$(SERVICE_NAME) IMAGE=${IMAGE} VERSION=${VERSION} && envsubst < deployment.yaml | kubectl apply -f -


undeploy: check-env
	export NAMESPACE=$(NAMESPACE) REALM=$(REALM) SERVICE_NAME=$(SERVICE_NAME) VERSION=${VERSION} && envsubst < deployment.yaml | kubectl delete -f -
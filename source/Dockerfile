FROM golang:1.22 as builder

WORKDIR /source

ARG TARGETOS TARGETARCH

ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH
ENV CGO_ENABLED 0 

ADD go.mod /source
ADD go.sum /source
WORKDIR /source
RUN go mod download
ENV GOCACHE /source/.cache

ADD . /source
RUN  --mount=type=cache,target=/source/.cache go build -ldflags="-w -s" -o istio-grpc-test .

# final image
#FROM docker.repo.splunkdev.net/observability/scratch:v0.4.0 AS final
FROM scratch as final

USER 10000
COPY --from=builder /source/server /usr/local/bin/istio-grpc-test

ENTRYPOINT ["istio-grpc-test"]

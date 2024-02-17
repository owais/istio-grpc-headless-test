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
RUN  --mount=type=cache,target=/source/.cache go build -ldflags="-w -s" -o bin/istio-grpc-test .

# final image
FROM scratch as final

USER 10000
COPY --from=builder /source/bin/istio-grpc-test /usr/local/bin/istio-grpc-test

ENTRYPOINT ["istio-grpc-test"]

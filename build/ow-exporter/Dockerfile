FROM registry.hub.docker.com/library/golang:1.17.6@sha256:2adfd0d7f507cfe0dab9aefaf3e0b973b9c8cce48caa43431dca7f8a2cef5557 AS build

RUN echo 'nobody:x:65534:65534:Nobody:/:' > /tmp/passwd \
    && apt-get update \
    && apt-get install -y --no-install-recommends upx-ucl=3.96-2

WORKDIR /go/src/github.com/lexfrei/tools/
COPY . /go/src/github.com/lexfrei/tools/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o ow-exporter ./cmd/ow-exporter/ow-exporter.go \
    && upx --best --lzma ow-exporter

FROM scratch

EXPOSE 9420/tcp

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/lexfrei/tools/ow-exporter /

USER nobody
ENTRYPOINT ["/ow-exporter"]

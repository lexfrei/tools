FROM registry.hub.docker.com/library/golang:1.20.1 AS build

RUN echo 'nobody:x:65534:65534:Nobody:/:' > /tmp/passwd \
    && apt-get update \
    && apt-get install -y --no-install-recommends upx-ucl=3.96-2

WORKDIR /go/src/github.com/lexfrei/tools/
COPY . /go/src/github.com/lexfrei/tools/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" ./cmd/me-site/me-site.go \
    && upx --best --lzma me-site
    

FROM scratch

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/lexfrei/tools/me-site /

EXPOSE 8080

USER nobody
ENTRYPOINT ["/me-site"]
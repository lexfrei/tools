FROM docker.io/library/golang:1.25.1-alpine3.22 AS build

RUN echo 'nobody:x:65534:65534:Nobody:/:' > /tmp/passwd && \
    apk add --no-cache upx=5.0.2-r0

WORKDIR /go/src/github.com/lexfrei/tools/
COPY . /go/src/github.com/lexfrei/tools/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" ./cmd/me-site/me-site.go && \
    upx --best --lzma me-site

FROM scratch

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/lexfrei/tools/me-site /

EXPOSE 8080

USER nobody
ENTRYPOINT ["/me-site"]
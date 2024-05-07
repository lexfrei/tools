FROM docker.io/library/golang:1.22.3-alpine3.18 AS build

RUN echo 'nobody:x:65534:65534:Nobody:/:' > /tmp/passwd && \
    apk add --no-cache upx=4.0.2-r0

WORKDIR /go/src/github.com/lexfrei/tools/
COPY . /go/src/github.com/lexfrei/tools/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o ow-exporter ./cmd/ow-exporter/ow-exporter.go && \
    upx --best --lzma ow-exporter

FROM scratch

EXPOSE 9420/tcp

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/lexfrei/tools/ow-exporter /

USER nobody
ENTRYPOINT ["/ow-exporter"]

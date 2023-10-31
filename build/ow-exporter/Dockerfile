FROM registry.hub.docker.com/library/golang:1.21.3 AS build

RUN echo 'nobody:x:65534:65534:Nobody:/:' > /tmp/passwd

WORKDIR /go/src/github.com/lexfrei/tools/
COPY . /go/src/github.com/lexfrei/tools/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o ow-exporter ./cmd/ow-exporter/ow-exporter.go

FROM scratch

EXPOSE 9420/tcp

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/lexfrei/tools/ow-exporter /

USER nobody
ENTRYPOINT ["/ow-exporter"]

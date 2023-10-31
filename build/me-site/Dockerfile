FROM registry.hub.docker.com/library/golang:1.21.3 AS build

RUN echo 'nobody:x:65534:65534:Nobody:/:' > /tmp/passwd

WORKDIR /go/src/github.com/lexfrei/tools/
COPY . /go/src/github.com/lexfrei/tools/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" ./cmd/me-site/me-site.go
    

FROM scratch

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/lexfrei/tools/me-site /

EXPOSE 8080

USER nobody
ENTRYPOINT ["/me-site"]
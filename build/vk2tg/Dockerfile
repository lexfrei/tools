FROM registry.hub.docker.com/library/golang:1.21.3 AS build

RUN echo 'nobody:x:65534:65534:Nobody:/:' > /tmp/passwd

WORKDIR /go/src/github.com/lexfrei/vk2tg/
COPY . /go/src/github.com/lexfrei/vk2tg/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" ./cmd/vk2tg/vk2tg.go

FROM scratch

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/lexfrei/vk2tg/vk2tg /

USER nobody
ENTRYPOINT ["/vk2tg"]

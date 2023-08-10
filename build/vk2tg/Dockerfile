FROM registry.hub.docker.com/library/golang:1.21.0 AS build

RUN echo 'nobody:x:65534:65534:Nobody:/:' > /tmp/passwd \
    && apt-get update \
    && apt-get install -y --no-install-recommends upx-ucl=3.96-2

WORKDIR /go/src/github.com/lexfrei/vk2tg/
COPY . /go/src/github.com/lexfrei/vk2tg/

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" ./cmd/vk2tg/vk2tg.go \
    && upx --best --lzma vk2tg

FROM scratch

COPY --from=build /tmp/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/github.com/lexfrei/vk2tg/vk2tg /

USER nobody
ENTRYPOINT ["/vk2tg"]

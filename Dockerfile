FROM golang:1.19.3-alpine AS builder
RUN apk add --no-cache git alpine-sdk
RUN go install -v github.com/voxelin/gh0.st@latest

FROM alpine:3.16.2
RUN apk -U upgrade --no-cache \
    && apk add --no-cache bind-tools ca-certificates chromium
COPY --from=builder /go/bin/gh0.st /usr/local/bin/

ENTRYPOINT ["gh0.st"]
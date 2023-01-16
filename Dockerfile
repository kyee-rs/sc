FROM golang:1.19.4-alpine AS builder
RUN apk add build-base
WORKDIR /go/src/ghost
COPY go.mod go.sum *.go /go/src/ghost/
RUN go build -o /go/bin/ghost

FROM alpine:3.17.0
COPY --from=builder /go/bin/ghost /usr/local/bin/

ENTRYPOINT ["ghost"]
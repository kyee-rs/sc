FROM golang:1.19-alpine AS builder
RUN apk add build-base
WORKDIR /code

ADD go.mod go.sum /code/
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go mod download -x

ADD . .
RUN --mount=type=cache,mode=0755,target=/go/pkg/mod go build -o /ghost .

FROM alpine
WORKDIR /
COPY --from=builder /ghost /usr/local/bin/
ENTRYPOINT ["ghost"]

FROM --platform=$BUILDPLATFORM golang:1.19-alpine AS builder
WORKDIR /code
ENV CGO_ENABLED=0

ADD go.mod go.sum /code/
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get

ADD . .
RUN go build -o /ghost .

FROM alpine
WORKDIR /
LABEL org.opencontainers.image.description "Ghost is a simple, fast, and reliable way to share your files."
COPY --from=builder /ghost /usr/local/bin/
ENTRYPOINT ["ghost"]
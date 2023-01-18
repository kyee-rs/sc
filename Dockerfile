FROM --platform=$BUILDPLATFORM golang:1.19-alpine AS builder
WORKDIR /code
ENV CGO_ENABLED=0

ADD go.mod go.sum /code/
RUN go mod download

ADD . .
RUN go build -o /ghost .

FROM alpine
WORKDIR /
COPY --from=builder /ghost /usr/local/bin/
ENTRYPOINT ["ghost"]
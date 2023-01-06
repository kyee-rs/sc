FROM golang:1.19-alpine
RUN apk add build-base
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /ghost
ENTRYPOINT [ "/ghost" ]
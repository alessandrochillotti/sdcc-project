FROM golang:1.16-alpine

WORKDIR /app

COPY ./node/node.go ./node/
COPY ./lib ./lib

RUN go mod init alessandro.it/app \
    && go mod tidy \
    && go mod download \
    && go build lib/lib.go

CMD [ "sh", "-c", "go run node/node.go" ]
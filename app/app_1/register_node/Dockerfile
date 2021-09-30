FROM golang:1.16-alpine

WORKDIR /app

COPY ./register_node/register_node.go ./register_node/
COPY ./lib ./lib

RUN go mod init alessandro.it/app \
    && go mod tidy \
    && go mod download \
    && go build lib/lib.go

CMD [ "sh", "-c", "go run register_node/register_node.go" ]
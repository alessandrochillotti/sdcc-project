FROM golang:1.16-alpine

ARG ALGO
ENV ALGORITHM ${ALGO}

WORKDIR /app

COPY ./node-algorithm-2/node.go ./node/node.go
COPY ./lib ./lib

RUN go mod init alessandro.it/app \
    && go mod tidy \
    && go mod download \
    && go build lib/lib.go \
    && go build -o main_node node/node.go

CMD [ "sh", "-c", "./main_node ${ALGORITHM}" ]
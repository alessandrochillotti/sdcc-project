FROM golang:1.16-alpine

ARG NUMBER_NODE
ENV NODES ${NUMBER_NODE}

WORKDIR /app

COPY ./register_node/ ./register_node/
COPY ./utils/structs.go ./utils/structs.go
COPY ./utils/utils.go ./utils/utils.go

RUN go mod init alessandro.it/app \
    && go mod tidy \
    && go mod download \
    && go build utils/structs.go

CMD [ "sh", "-c", "go run register_node/register_node.go" ]
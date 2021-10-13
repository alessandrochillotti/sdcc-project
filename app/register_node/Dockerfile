FROM golang:1.16-alpine

ARG NUMBER_NODE
ENV NODES ${NUMBER_NODE}

WORKDIR /app

COPY ./register_node/ ./register_node/
COPY ./utils/structs.go ./utils/structs.go
COPY ./utils/utils.go ./utils/utils.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

CMD [ "sh", "-c", "go run register_node/register_node.go" ]
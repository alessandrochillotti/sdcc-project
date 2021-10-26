FROM golang:1.16-alpine

ARG NUMBER_NODE
ENV NODES ${NUMBER_NODE}

WORKDIR /app

COPY ./sequencer/ ./sequencer/
COPY ./utils/structs.go ./utils/structs.go
COPY ./utils/utils.go ./utils/utils.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

CMD [ "sh", "-c", "go run sequencer/sequencer.go" ]
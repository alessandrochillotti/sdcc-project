FROM golang:1.16-alpine

ARG ALGO
ARG NUMBER_NODE
ENV ALGORITHM ${ALGO}
ENV NODES ${NUMBER_NODE}

WORKDIR /app

COPY ./peer/ ./peer/
COPY ./utils/ ./utils/
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

CMD [ "sh", "-c", "go run peer/*.go" ]
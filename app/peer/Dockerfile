FROM golang:1.16-alpine

ARG ALGO
ARG NUMBER_NODE
ENV ALGORITHM ${ALGO}
ENV NODES ${NUMBER_NODE}

WORKDIR /app

COPY ./peer/ ./peer/
COPY ./utils/ ./utils/

RUN go mod init alessandro.it/app \
    && go mod tidy \
    && go mod download 

CMD [ "sh", "-c", "go run peer/*.go" ]
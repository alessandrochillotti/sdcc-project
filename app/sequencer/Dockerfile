FROM golang:1.16-alpine

WORKDIR /app

COPY ./sequencer/ ./sequencer/
COPY ./utils/structs.go ./utils/structs.go

RUN go mod init alessandro.it/app \
    && go mod tidy \
    && go mod download \
    && go build utils/structs.go

CMD [ "sh", "-c", "go run sequencer/sequencer.go" ]
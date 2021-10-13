FROM golang:1.16-alpine

WORKDIR /app

COPY ./sequencer/ ./sequencer/
COPY ./utils/structs.go ./utils/structs.go
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

CMD [ "sh", "-c", "go run sequencer/sequencer.go" ]
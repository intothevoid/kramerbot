FROM golang:latest

RUN go build .

ENTRYPOINT [ "./kramerbot" ]

FROM golang:alpine

# Set the working directory
WORKDIR /kramerbot

# Copy the source code
ADD . /kramerbot

# Build the bot
RUN CGO_ENABLED=0 go build .

ENTRYPOINT [ "./kramerbot" ]

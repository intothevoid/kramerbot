FROM golang:latest

# Set the working directory
WORKDIR /kramerbot

# Copy the source code
ADD . /kramerbot

# Build the bot
RUN go build .

ENTRYPOINT [ "./kramerbot" ]

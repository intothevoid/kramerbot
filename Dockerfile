FROM golang:latest

# Set the working directory
WORKDIR /kramerbot

# Copy the source code
ADD . /kramerbot

# Create directories for data and webapp
RUN mkdir -p /kramerbot/data /kramerbot/webapp/static

# Build the bot
RUN go build .

# Expose ports for the webapp
EXPOSE 5173

ENTRYPOINT [ "./kramerbot" ]

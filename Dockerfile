FROM golang:latest

# Install sqlite3
RUN apt-get update && apt-get install -y sqlite3

WORKDIR /app

# Create data directory with correct permissions
RUN mkdir -p /app/data && chmod 777 /app/data

# Copy source code
ADD . /app

# Build the bot
RUN go build .

# Set the entry point
ENTRYPOINT ["./kramerbot"]

# ── Stage 1: Build React frontend ────────────────────────────────────────────
FROM node:20-slim AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ── Stage 2: Build Go binary ──────────────────────────────────────────────────
FROM golang:1.24 AS go-builder
WORKDIR /app

# Install gcc (required by mattn/go-sqlite3)
RUN apt-get update && apt-get install -y gcc

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Embed the built React app so the Go binary can serve it
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=1 GOOS=linux go build -o kramerbot .

# ── Stage 3: Minimal runtime image ───────────────────────────────────────────
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*

WORKDIR /app
RUN mkdir -p /app/data && chmod 777 /app/data

COPY --from=go-builder /app/kramerbot ./kramerbot
COPY config.yaml ./config.yaml

# API port
EXPOSE 8080

ENTRYPOINT ["./kramerbot"]

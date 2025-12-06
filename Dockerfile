# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o mss-bot ./cmd/bot

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/mss-bot .

# Copy config example
COPY --from=builder /app/configs ./configs

# Create data directory
RUN mkdir -p /app/data

ENTRYPOINT ["./mss-bot"]
CMD ["-config", "/app/configs/config.kdl"]

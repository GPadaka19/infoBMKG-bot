# Stage 1: Builder
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates (needed for go get and HTTPS)
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary (better for alpine/scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -o bmkg-bot cmd/main.go

# Stage 2: Runner
FROM alpine:latest

# Install CA certificates for HTTPS connections to BMKG/Telegram
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bmkg-bot .

# Create a non-root user for security
RUN adduser -D botuser
USER botuser

# We don't need to EXPOSE any port because this is a polling bot (outbound client)

# Run the binary
CMD ["./bmkg-bot"]

# Use the official Go image as base
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum .env.dev ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN GIN_MODE=release go build -o ./main .

# Use a minimal base image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and sqlite
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create .env file placeholder (will be overridden by environment variables)
RUN touch .env.dev

# Expose port
EXPOSE 4000

# Command to run
CMD ["./main"]

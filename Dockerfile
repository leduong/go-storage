# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o app .

# Production stage
FROM alpine:latest
WORKDIR /app

# Copy built binary and necessary files
COPY --from=builder /app/app .

# Install required runtime packages
RUN apk add --no-cache ca-certificates

# Expose port and run the application
EXPOSE 8080
CMD ["./app"]

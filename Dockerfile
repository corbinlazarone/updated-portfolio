# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy everything
COPY . .

# Build the application - compile all files in cmd package
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/*.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy UI directory for templates/static files
COPY --from=builder /app/ui ./ui

# Render provides the PORT environment variable
EXPOSE ${PORT:-8080}

# Run the application
CMD ["./main"]

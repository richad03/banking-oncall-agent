# Multi-stage build for smaller final image
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o banking-agent cmd/bot/main.go

# Final stage - minimal runtime image
FROM alpine:3.18

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' banking-agent

# Set working directory
WORKDIR /home/banking-agent

# Copy binary from builder stage
COPY --from=builder /app/banking-agent .

# Copy knowledge base (if needed in container)
# COPY --from=builder /app/banking-oncall ./banking-oncall

# Change ownership to non-root user
RUN chown -R banking-agent:banking-agent /home/banking-agent

# Switch to non-root user
USER banking-agent

# Expose port (if using HTTP health checks)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep banking-agent || exit 1

# Run the application
CMD ["./banking-agent"]
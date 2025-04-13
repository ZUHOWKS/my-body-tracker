# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the project
COPY api api
COPY internal internal
COPY cmd cmd

# Build the application
RUN go build -o bodytracker_api ./cmd/api

# Final stage
FROM scratch

# Copy the binary from builder
COPY --from=builder /app/bodytracker_api /app/bodytracker_api

# Expose portz
EXPOSE 8080

# Run the application
CMD ["/app/bodytracker_api"]

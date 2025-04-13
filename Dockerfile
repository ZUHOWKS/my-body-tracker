# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bodytracker_api ./api

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/bodytracker_api .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./bodytracker_api"]

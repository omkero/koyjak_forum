# Start from the official Golang image as a build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
RUN go build -o app cmd/server/main.go
# Use a multi-stage build to keep the final image small

# Use a minimal image for the final container
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .
COPY views ./views
COPY .env .env
COPY public ./public
COPY config ./config

# Set environment variables if needed
# ENV DATABASE_URL=...

# Expose port (change if your app uses a different port)
EXPOSE 8080

# Run the binary
CMD ["./app"]
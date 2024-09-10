# Use the official Golang image as the build environment
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/taskmanager

# Start a new stage from scratch
FROM alpine:latest

# Install necessary packages (including ca-certificates for HTTPS)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the .env file into the container
COPY .env .env

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]
# Use the official Golang image as a build environment
FROM golang:1.22.5 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifest and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
WORKDIR /app/cmd/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# Use a minimal base image
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/cmd/user-service/main .

# Ensure the binary has execute permissions
RUN chmod +x ./main

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./main"]

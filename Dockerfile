# Start from the latest golang base image
FROM golang:1.21 AS builder

# Add Maintainer Info
LABEL maintainer="Russell Jones <jonesrussell42@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /app/main /app/main

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./app/main"]
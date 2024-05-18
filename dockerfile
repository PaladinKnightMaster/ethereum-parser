# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder

# Install Git; required for fetching the dependencies.
RUN apk add --no-cache git

# Setting up environment variables
ENV GOPROXY=direct

# Create and change to the app directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy only the necessary source files
COPY cmd/ ./cmd/
COPY config/ ./config/
COPY internal/api ./internal/api
COPY internal/parser ./internal/parser
COPY shared/ ./shared/

# Build the Go app
RUN go build -o main ./cmd/main.go

# Stage 2: Run the Go application
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

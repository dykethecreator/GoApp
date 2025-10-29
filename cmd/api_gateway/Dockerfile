# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod, go.sum and go.work to download dependencies
COPY go.mod go.sum go.work ./
RUN go mod download
RUN go work sync

# Copy the rest of the source code
COPY . .

# Build the api_gateway binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api_gateway ./cmd/api_gateway

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /bin/api_gateway /bin/api_gateway

# Expose the HTTP port
EXPOSE 8080

# Command to run the application
CMD ["/bin/api_gateway"]

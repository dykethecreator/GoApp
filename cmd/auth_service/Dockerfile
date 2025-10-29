# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Workspace modunu kapat: Docker içinde tek modül gibi davransın
ENV GOWORK=off

# Copy the rest of the source code
COPY go.mod go.sum ./

# Tüm bağımlılıkları projenin içindeki vendor klasöründen almasını söyle.
# Bu, internet erişimi olmadan derleme yapmamızı sağlar.
RUN go mod download

# Şimdi projenin geri kalanını ve en önemlisi vendor klasörünü kopyala
COPY . .
# Build the auth_service binary
# CGO_ENABLED=0 is important for creating a static binary that runs in a minimal container
# -o /bin/auth_service specifies the output path for the compiled binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth_service ./cmd/auth_service

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /bin/auth_service /bin/auth_service

# Expose the gRPC port that the service will listen on
# This is for documentation; the actual port mapping is in docker-compose.yml
EXPOSE 50051

# Command to run the application
CMD ["/bin/auth_service"]

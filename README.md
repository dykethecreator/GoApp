# WhatsApp Clone Backend

This project is a Go-based backend for a WhatsApp-like application, built with a microservices architecture.

## Project Structure

The project follows a clean architecture and microservices pattern.

- `/cmd`: Entry points for each service (`main.go`).
- `/internal`: Private application and business logic for each service.
- `/pkg`: Shared libraries and domain types used across services.
- `/proto`: gRPC protocol definitions for inter-service communication.
- `/migrations`: Database schema migrations.
- `/configs`: Configuration files for different environments.

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.18 or higher

### Running the application

1.  **Start the infrastructure:**
    ```bash
    docker-compose up -d
    ```

2.  **Run database migrations:**
    You'll need a migration tool like `golang-migrate/migrate`.
    ```bash
    migrate -database "postgres://postgres:Fbtex1967.@localhost:5432/whatsapp_clone_dev?sslmode=disable" -path migrations up
    ```

3.  **Run the services:**
    Navigate to each service's directory and run it.
    ```bash
    go run ./cmd/api_gateway/
    go run ./cmd/auth_service/
    # ... and so on for other services
    ```

## Quickstart: AuthService + Postman gRPC (Windows)

End-to-end minimum setup to test OTP and JWT issuance via Postman using gRPC.

### 1) Environment variables

This service auto-loads environment files based on context:

- Local runs: loads `.env.local` (we added one with sensible defaults)
- Docker runs: `docker-compose` passes `.env.docker` and also sets `RUNNING_IN_DOCKER=true` (the app attempts to load `.env.docker` but works even if the file isn't baked into the image)
- Base `.env` is also loaded last if present (for overrides)

Files created for you:

- `.env.local` — connects to Postgres at `localhost:5432` and sets JWT/Twilio placeholders
- `.env.docker` — connects to Postgres at `postgres:5432` (Compose service name) and sets the same placeholders

Edit these files as needed:

- Set `JWT_SECRET` to a strong value
- Add real Twilio credentials if you want to use OTP verification against Twilio

### 2) Start PostgreSQL (Docker)

Run only Postgres in the background:

```powershell
docker-compose up -d postgres
```

### 3) Apply database schema (migrations)

Apply the initial schema directly via psql inside the Postgres container (TCP to avoid socket issues):

```powershell
$pg = docker ps --filter "name=postgres" --format "{{.ID}}"

# Ensure UUID extension exists
docker exec -i $pg psql -h 127.0.0.1 -U user -d whatsapp_clone_dev -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

# Apply the schema
Get-Content -Raw .\migrations\0001_initial_schema.up.sql | docker exec -i $pg sh -c 'psql -h 127.0.0.1 -U user -d whatsapp_clone_dev -v ON_ERROR_STOP=1'
```

If you previously ran into an inline index syntax error for `call_logs`, this repository already fixes it by creating the index separately.

### 4) Run Auth Service locally

From the project root:

```powershell
go run ./cmd/auth_service
```

You should see a log similar to: `auth_service listening on :50051 (env=local)`.

### 5) Test via Postman (gRPC)

1. In Postman, create a new gRPC Request.
2. Server URL: `localhost:50051` (plaintext/TLS off).
3. Import `proto/auth.proto`.
4. Select `auth.AuthService` and call:
     - `SendOTP` with body `{ "phone_number": "+9055xxxxxxx" }`
     - `VerifyOTP` with body `{ "phone_number": "+9055xxxxxxx", "otp_code": "123456" }`

On success, `VerifyOTP` returns `access_token` and `refresh_token`.

5a) Token utilities via gRPC

- `ValidateToken` with body `{ "access_token": "<ACCESS>" }` → returns `{ is_valid, user_id }` (no error for invalid; just `is_valid=false`).
- `RefreshToken` with body `{ "refresh_token": "<REFRESH>" }` → returns `{ access_token, refresh_token }` (refresh token rotasyonu etkin).

5b) Revoke sessions via gRPC

- `RevokeCurrentDevice` with body `{ "refresh_token": "<REFRESH>" }` → returns `{ success: true }`. Afterwards, the same refresh token can no longer be used.
- `LogoutAllDevices` with body `{ "access_token": "<ACCESS>" }` → returns `{ success: true }`. Afterwards, any existing refresh tokens for that user are invalidated (server checks DB-stored hashes and sees they are revoked).

## Notes: Local vs Docker run

- Local app run (recommended for quick testing):
    - App connects to Dockerized Postgres via `localhost:5432`.
    - `.env.local` already contains `DATABASE_URL=...@localhost:5432/...`.
    - Start with `go run ./cmd/auth_service`.

- Docker app run (Compose):
    - App connects via the Compose network using host `postgres`.
    - `docker-compose up --build` will build and run `auth_service` against `postgres`.
    - `docker-compose` passes `.env.docker` to the container.

## Troubleshooting

- `unknown driver "postgres"`: Make sure the project has `github.com/lib/pq` and the driver is blank-imported (already added in `pkg/database/postgres.go`).
- `relation "..." does not exist`: Ensure you applied migrations to the exact database your service is using.
- `function uuid_generate_v4() does not exist`: Run `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";` on the target database, or switch to `pgcrypto` + `gen_random_uuid()` in schema.

## Services

- **API Gateway**: The single entry point for all client requests.
- **Auth Service**: Handles user authentication, registration, and session management.
- **Chat Service**: Manages chat rooms and user memberships.
- **Message Worker**: Asynchronously processes and stores messages from a queue.
- **Realtime Service**: Manages WebSocket connections for real-time communication.
- **Status Service**: Handles user status updates (stories).

- Revocation and middleware (overview)
    - Refresh tokens: We added a migration (`0002_user_devices_revocation.up.sql`) to support storing refresh token hashes and revocation timestamps in `user_devices`.
    - Middleware: `internal/auth/middleware/auth_interceptor.go` bir gRPC unary interceptor sağlar; `authorization: Bearer <token>` başlığından access token doğrular, `user_id`’yi context’e ekler.
    - AuthService içinde interceptor varsayılan olarak AuthService RPC’lerini muaf tutar (public uçlar). Diğer servislerde enable etmek için:
        - Aynı `JWT_SECRET` ile bir `TokenManager` oluşturun (örn. 15dk access / 7gün refresh).
        - gRPC server’ı `grpc.UnaryInterceptor(UnaryAuthInterceptor(tm))` ile oluşturun.
        - Gerekirse `info.FullMethod` ile public uçları muaf tutabilirsiniz.

    - Applying the revocation migration:
        ```powershell
        # Local Postgres
        psql -U postgres -h localhost -d whatsapp_clone_dev -f .\migrations\0002_user_devices_revocation.up.sql
        ```
    - Next steps to fully wire revocation:
        - On VerifyOTP: hash the refresh token (e.g., SHA-256), store in `user_devices` with `last_login_at=NOW()`.
        - On RefreshToken: compute the hash and ensure an active (revoked_at IS NULL) record exists for that user and hash; otherwise return Unauthenticated.
        - Optionally rotate refresh token and update storage; return the new refresh token in the response (requires proto change).

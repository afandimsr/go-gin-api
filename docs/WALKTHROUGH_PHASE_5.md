# Phase 5: Client Auth Integration Walkthrough

I have implemented the integration with an external Client Auth Service.

## Changes

### 1. Configuration
- Added `CLIENT_AUTH_URL` to `.env` and `internal/config/config.go`.

### 2. External Auth Service
- Created `internal/infrastructure/external/auth_service.go` implementing `AuthClient`.
- Defines `Login(email, password string) (bool, error)` to POST credentials to the external URL.

### 3. User UseCase
- Updated `UserUsecase` to accept `AuthService` interface.
- Modified `Login` logic:
    1.  Check `AuthService` first (if configured).
    2.  If external auth succeeds, login is successful.
    3.  If external auth fails or is not configured, fallback to local `bcrypt` check.

### 4. Bootstrap
- Updated `internal/bootstrap/app.go` to initialize `AuthClient` with the configured URL and pass it to the UseCase.

## How to Verify

1.  **Configure External Auth**:
    Set `CLIENT_AUTH_URL=http://external-auth-service.com` in `.env`.

2.  **Test Login**:
    call `POST /api/v1/login`. The system will attempt to validate against the external URL first.

3.  **Run Tests**:
    ```bash
    go test -v ./internal/usecase/user/...
    ```

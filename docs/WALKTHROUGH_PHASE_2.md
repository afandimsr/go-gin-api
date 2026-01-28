# Phase 2: Authentication & Security Walkthrough

I have completed the implementation of authentication using JWT and bcrypt.

## Changes

### 1. Database & Entity
- Added `password` column to `users` table via migration: `20260106150500_add_password_to_users.up.sql`.
- Updated `User` entity to include `Password` field.
- Updated Repository to handle password storage.

### 2. Password Hashing
- Integrated `bcrypt` in `UserUsecase` for `Create` and `Update` methods.

### 3. Login Feature
- Added `Login` method in `UserUsecase`.
- Added `Login` handler in `UserHandler`.
- Endpoint: `POST /api/v1/login`
- Returns a JWT token on success.

### 4. JWT Middleware
- Created `internal/delivery/http/middleware/auth_middleware.go`.
- Validates `Authorization: Bearer <token>` header.
- Protects `UPDATE` and `DELETE` routes for users.

## How to Test

1.  **Register**: `POST /api/v1/users` with `name`, `email`, `password`.
2.  **Login**: `POST /api/v1/login` with `email`, `password`. Copy the `token`.
3.  **Get User**: `GET /api/v1/users/:id` (Public).
4.  **Update User**: `PUT /api/v1/users/:id` (Protected). Add Header `Authorization: Bearer <token>`.

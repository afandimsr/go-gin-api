# Phase 1: User CRUD Walkthrough

I have completed the implementation of the full CRUD lifecycle for the User Domain.

## Changes

### 1. Repository
[user_repository.go](file:///d:/Project/Backend/Go/Framework/Go%20Gin/go-gin-api/internal/infrastructure/persistent/mysql/repository/user_repository.go)
- Added `FindByID`, `Update`, `Delete`.

### 2. Usecase
[user_usecase.go](file:///d:/Project/Backend/Go/Framework/Go%20Gin/go-gin-api/internal/usecase/user/user_usecase.go)
- Implemented business logic for CRUD.

### 3. Handlers & Routes
- GET `/api/v1/users`
- GET `/api/v1/users/:id`
- POST `/api/v1/users`
- PUT `/api/v1/users/:id`
- DELETE `/api/v1/users/:id`

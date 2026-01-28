# Project Plan: Go Gin API

This document serves as the implementation plan, including the Product Requirement Document (PRD) and Minimum Viable Product (MVP) definition for the **Go Gin API** project.

## Product Requirement Document (PRD)

### 1. Project Overview
**Name**: Go Gin API (User Management System)
**Description**: A RESTful API built with Go and the Gin framework, following Clean Architecture principles. It serves as a backend foundation for managing users and potentially other entities, providing secure and scalable data access.

### 2. Goals
- **Maintainability**: Follow Clean Architecture to separate concerns (Domain, Usecase, Repository, Delivery).
- **Scalability**: Designed to be easily extended with new modules (e.g., Auth, Products).
- **Documentation**: Fully documented API using Swagger.
- **Reliability**: Robust error handling and database interactions.

### 3. User Stories (MVP)
- **User Management**:
    - As a user, I want to create an account (Register).
    - As a user, I want to retrieve my profile details.
    - As an admin, I want to see a list of all users.
    - As a user, I want to update my profile.
    - As a user, I want to delete my account.
- **Authentication**:
     - As a user, I want to login with my email and password.
     - As a user, I want my sensitive actions (update/delete) to be protected.
     - **[NEW]** As a system admin, I want to configure an external Client Auth Service to validate credentials instead of the local database.

### 4. Non-Functional Requirements
- **Performance**: Response time under 200ms for typical requests.
- **Database**: MySQL for persistent storage.
- **Security**: Basic input validation and potential future JWT authentication.
- **Code Quality**: Go linting and idiomatic Clean Architecture structure.

---

## Minimum Viable Product (MVP) Definition

The MVP focuses on a complete User Management cycle.

### Core Features
1.  **User CRUD**:
    - `POST /users`: Create a new user.
    - `GET /users`: Retrieve all users (Admin view).
    - `GET /users/:id`: Retrieve specific user.
    - `PUT /users/:id`: Update user details.
    - `DELETE /users/:id`: Delete a user.
2.  **Authentication**:
    - `POST /login`: Authenticate and receive JWT.
    - Middleware to protect routes.
    - **External Auth Support**: Optional integration with a Client Auth Service.

3.  **Infrastructure**:
    - MySQL Database connection.
    - Database Migrations (User table).
    - Swagger Documentation (`/swagger/index.html`).
    - Global Error Handling Middleware.

---

## Roadmap & Implementation Plan

### Phase 1: Foundation & User Domain (Completed)
- [x] Project Structure Setup (Clean Architecture).
- [x] Database Connection (MySQL).
- [x] User Entity & Repository Interface.
- [x] User Repository Implementation (MySQL).
- [x] Migration Setup (Users Table).
- [x] User Usecase Logic (CRUD).
- [x] HTTP Handlers & Routing.
- [x] Request Validation.

### Phase 2: Authentication & Security (Completed)
- [x] Add `password` column to users.
- [x] Implement Password Hashing (bcrypt).
- [x] Implement Login Endpoint (`POST /login`).
- [x] Implement JWT Middleware.

### Phase 3: Enhancements (Completed)
- [x] **Pagination**:
    - Update `FindAll` to accept `limit` and `offset`.
    - Update `GetAll` usecase to accept `page` and `limit`.
    - Update `GetUsers` handler to parse query params.
- [x] **Unit Tests**:
    - Use `testify/mock` to mock `UserRepository`.
    - Test `Create`, `GetByID`, `Login` methods in `UserUsecase`.
- [x] **CI/CD**:
    - Create `.github/workflows/go.yml` to run tests and build.

### Phase 4: Developer Experience (Completed)
- [x] Add Makefile / Scripts.

### Phase 5: Client Auth Integration (Completed)
- [x] **Configuration**: `CLIENT_AUTH_URL`.
- [x] **Infrastructure**: `AuthClient` implementation.
- [x] **UseCase**: Update `Login` to check external service first.

## Output Documentation
- [Phase 1 Walkthrough](docs/WALKTHROUGH_PHASE_1.md)
- [Phase 2 Walkthrough](docs/WALKTHROUGH_PHASE_2.md)
- [Phase 3 Walkthrough](docs/WALKTHROUGH_PHASE_3.md)
- [Phase 5 Walkthrough](docs/WALKTHROUGH_PHASE_5.md)

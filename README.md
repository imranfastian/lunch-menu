# NBIS Assignment

Assignment for the position as backend system developer with focus on security  
Reference number: UFV-PA 2025/2222  
Submitted by: Imran Bashir

---

## Lunch Menu API - Secure Restaurant & Menu Management

A robust REST API built with Go and PostgreSQL for managing restaurants and their menus, with a strong focus on security, authentication, authorization, and consistent API responses.

---

## Getting Started

### 1. Clone the Repository

```bash
git clone <repository-url>
cd lunch_menu
```

### 2. Environment Variables

Make sure to have a `.env` file in the root directory with the configuration settings. You can copy from `.env` and modify as needed.

### 3. Install Dependencies

```bash
go mod tidy
```

### 4. Generate Swagger Documentation

```bash
swag init --parseDependency --parseInternal

```

### 5. Run with Docker Compose

```bash
docker compose --env-file .env up --build
```

### 6. Access the API

The API will be running at `http://localhost:8000`.
Swagger documentation is available at `http://localhost:8000/swagger/index.html#/`

### 7. Run Tests

```bash
go test  ./internal/tests/...
go test -v ./internal/tests/...
```

---

## Features

- **Consistent Unified API Responses:**  
  All endpoints return a standard JSON structure using `models.StandardResponse`, ensuring predictable and easy-to-parse responses for both success and error cases.

- **Authentication & Authorization:**

  - **User Registration:** Admins can register new users via `/api/user/register`.
  - **Login:** Users log in via `/api/user/login` and receive an access token (JWT) and a refresh token.
  - **Access Token:** JWT is generated based on user information and token expiry.
  - **Refresh Token:** Used to renew access tokens without re-authentication.
  - **Logout:** `/api/user/logout` endpoint resets cookies and blacklists the access token, preventing its further use.
  - **Token Blacklisting:** After logout, the access token is blacklisted, so it cannot be used for any further requests.
  - **Middleware:**
    - **Auth Middleware:** Protects endpoints requiring authentication, checks for valid and non-blacklisted JWTs.
    - **Role-based Access:** Only admins can perform create, update, or delete operations on restaurants and menu items.
    - **CSRF Protection:** CSRF tokens are set and validated for state-changing operations.
    - **CORS:** Configured for secure cross-origin requests.
    - **Rate Limiting:** Configurable middleware for API abuse protection.

- **Public and Protected Endpoints:**

  - **Public:**
    - List restaurants and menu items
    - Get restaurant or menu item details
  - **Protected (Authentication Required):**
    - Create, update, delete restaurants and menu items
    - User registration, login, logout

- **JWT & Cookie Management:**

  - Access and refresh tokens are set as HTTP-only cookies.
  - Secure cookie handling for authentication and session management.

- **Automatic Database Migrations:**

  - GORM AutoMigrate ensures the database schema matches the models on startup.

- **Swagger/OpenAPI Documentation:**
  - Interactive API docs available at [http://localhost:8000/swagger/index.html#/](http://localhost:8000/swagger/index.html#/)

---

## API Response Structure

All API responses use a unified format for consistency and ease of use. models.StandardResponse for both success and error cases.

```json
{
  "message": "Description of the result",
  "data": {},
  "error": {
    "error": "ERROR_CODE",
    "message": "Detailed error message"
  }
}
```

_`data` is present on success; `error` is present on error._

---

## Example Endpoints

### Public

- `GET /api/restaurants` — List restaurants
- `GET /api/restaurants/{id}` — Get restaurant details
- `GET /api/restaurants/{id}/menu` — List menu items for a restaurant
- `GET /api/menu-items/{id}` — Get menu item details

### Protected (Admin Only)

- `POST /api/restaurants` — Create restaurant
- `PUT /api/restaurants/{id}` — Update restaurant
- `DELETE /api/restaurants/{id}` — Soft delete restaurant (`is_active=false`)
- `POST /api/menu-items` — Create menu item
- `PUT /api/menu-items/{id}` — Update menu item
- `DELETE /api/menu-items/{id}` — Soft delete menu item (`is_available=false`)

### User Management

- `POST /api/user/register` — Register admin user
- `POST /api/user/login` — Login and receive tokens and cookies
- `POST /api/user/logout` — Logout and blacklist token

### Statistics

- `GET /api/stats` — Returns business analytics and statistics

---

## Authentication & Security Flow

### Register User

- `POST /api/user/register`  
  Registers a new user with role and uniq username and email.
  newly registered users have `is_active=false` by default and cannot login until activated.
  newly registered users have role but cannot login other then admin role.

### Login

- `POST /api/user/login`  
  Registered user need to set true is_active to login.
  Registered user with admin role can login only at this moment.
  Returns access and refresh tokens as HTTP-only cookies.  
  Access token is a JWT containing user info and expiry.  
  Refresh token allows silent renewal of access token.

## Deleted restaurant or menu item

- Soft delete by setting `is_active=false` for restaurants and `is_available=false` for menu items.

### Access Protected Endpoints

- Send requests with the access token cookie (or Authorization header).
- Middleware checks token validity and blacklist status.

### Token Renewal

- If the access token expires, use the refresh token to obtain a new access token.

### Logout

- `POST /api/user/logout`  
  Clears authentication cookies and blacklists the access token to prevent further use.
  Clears refresh token from database.

---

## User Registration: IsActive Behavior

**Update:**  
Previously, newly registered users had `is_active=false` by default and required manual activation in the database before they could log in.  
**Now,** newly registered users have `is_active=true` by default, allowing immediate login and easier testing of authentication and other features.

- This change streamlines development and testing.
- For production, you may want to revert to manual activation for user approval workflows.

---

## Task 1: Secure Authentication & Authorization

**Task 1:**  
Implemented access control mechanisms to restrict administrative endpoints to authorized users only.

- Public endpoints remain open for menu and restaurant viewing.
- Admin endpoints (create, update, delete) require authentication and admin role.
- JWT-based authentication, secure cookie handling, and role-based middleware are used.
- CSRF protection, CORS, and token blacklisting are implemented for enhanced security.

---

## Task 2: Secure CRUD Functionality for Administrators

**Task 2:**  
Developed secure CRUD endpoints for restaurant and menu management, accessible only to authenticated admins.

- Admins can create, update, and soft-delete restaurants and menu items.
- All changes are protected by authentication and role-based authorization middleware.
- Consistent API response structure and error handling.
- Soft delete is implemented by toggling `is_active` or `is_available` flags.

---

## Optional Task 3: Kubernetes Deployment & Production Readiness

**Task 3:**  
A complete Kubernetes deployment setup is provided for production-ready deployment, configuration management, and security best practices.

- All Kubernetes manifests and configuration files are located in the [`k8s/`](./k8s/) directory.
- Includes:
  - API and PostgreSQL deployments
  - ConfigMap and Secret management
  - Ingress controller (NGINX) setup for load balancing and external routing
  - HTTPS/TLS configuration with self-signed certificates for secure access
  - Persistent storage for database
  - Security recommendations for production

**See [`k8s/README.md`](./k8s/README.md) for full instructions and details.**

---

## Project Directory Structure

```
lunch_menu/
├── main.go
├── init_db.sql
├── internal/
│   ├── handlers/      # HTTP handlers (controllers)
│   ├── models/        # GORM models and DTOs
│   ├── database/      # DB connection and CRUD
│   ├── middleware/    # Auth, CSRF, rate limiting
│   ├── routes/        # Route definitions
│   └── utils/         # JWT, response helpers, etc.
├── docs/              # Swagger docs (docs.go, swagger.json, swagger.yaml)
├── k8s/               # Kubernetes manifests, including secrets management, ingress, etc.
├── .env               # Environment variables (not committed)
└── ...
```

---

## Security Best Practices

- Never store secrets in images or code.
- Set COOKIE_SECURE=true and COOKIE_HTTPONLY=true in production.
- Restrict CORS to only trusted frontend domains.
- Use resource requests/limits for containers.
- Enable logging and monitoring (e.g., with Prometheus, Grafana, Loki).
- Use NetworkPolicies to restrict traffic between pods.
- Enable automatic restarts and health checks (liveness/readiness probes).
- Keep your images up to date and scan for vulnerabilities.

---

## Technologies Used

- Go (Golang)
- Gin Web Framework
- GORM (ORM for Go)
- PostgreSQL
- Swagger (API Documentation)
- Docker & Docker Compose

---

## Author

Imran Bashir  
Email: imranfastian@gmail.com

---

## Completion & Submission

All required tasks for the NBIS Assignment (reference number: UFV-PA 2025/2222) have been completed:

- **Task 1:** Secure authentication and authorization for admin endpoints.
- **Task 2:** Secure CRUD endpoints for administrators.
- **Optional Task 3:** Kubernetes deployment setup with configuration management and production security considerations (see [`k8s/README.md`](./k8s/README.md)).

The solution is documented and can be deployed by following the instructions in this repository.

Thank you for reviewing my submission!

---

# Procurement System - Backend

This directory contains the Go backend for the Procurement System.

## Setup

1.  **Database:**
    *   Make sure you have PostgreSQL running.
    *   Create a database (e.g., `procurement`).
    *   Connect to the database and run the schema definition in `migrations/001_initial_schema.sql` to create the necessary tables.

2.  **Environment Variables:**
    *   Copy the `.env.example` file to `.env`.
    *   Update the `.env` file with your database connection string and a secure JWT secret.
    ```
    DATABASE_URL=postgres://youruser:yourpassword@localhost:5432/procurement?sslmode=disable
    JWT_SECRET=a-very-secure-secret-key
    PORT=8080
    ```

3.  **Run the Server:**
    *   Navigate to the `backend` directory.
    *   Install dependencies: `go mod tidy`
    *   Run the server: `go run ./cmd/main.go`
    *   The server will start on the port specified in your `.env` file (defaults to 8080).

## Database Seeding

To populate the database with sample data for development and testing, you can run the seeder script. This will clean all existing data and create a set of users, vendors, and requisitions.

```bash
go run ./cmd/seeder/main.go
```

## API Endpoints

All endpoints are prefixed with `/api`.

### Authentication

*   **`POST /register`**
    *   **Description:** Registers a new user.
    *   **Body:**
        ```json
        {
          "name": "Test User",
          "email": "test@example.com",
          "password": "password123",
          "role": "Employee"
        }
        ```
    *   **Response:** `201 Created` with user object (without password).

*   **`POST /login`**
    *   **Description:** Authenticates a user and returns a JWT.
    *   **Body:**
        ```json
        {
          "email": "test@example.com",
          "password": "password123"
        }
        ```
    *   **Response:** `200 OK` with JWT token.
        ```json
        {
          "token": "your.jwt.token"
        }
        ```

### Profile Management

*All profile routes require authentication.*

*   **`GET /profile/me`**: Returns the profile of the currently logged-in user.
*   **`PUT /profile/me`**: Updates the logged-in user's name.
*   **`PUT /profile/password`**: Changes the logged-in user's password.

### User Management (Admin Only)

*All user management routes require a valid JWT from an "Admin" user.*

*   **`GET /users`**: Returns a list of all users.
*   **`GET /users/{id}`**: Returns a single user by ID.
*   **`PUT /users/{id}`**: Updates a user's name and role.
*   **`DELETE /users/{id}`**: Deletes a user.

### Vendor Management (Admin Only)

*All vendor routes require a valid JWT from an "Admin" user.*

*   **`POST /vendors`**: Creates a new vendor.
*   **`GET /vendors`**: Returns a list of all vendors.
*   **`GET /vendors/{id}`**: Returns a single vendor by ID.
*   **`PUT /vendors/{id}`**: Updates a vendor's details.
*   **`DELETE /vendors/{id}`**: Deletes a vendor.

### Purchase Requisitions

*All requisition routes require authentication.*

*   **`POST /requisitions`**
    *   **Description:** Creates a new purchase requisition. `requester_id` is taken from the JWT.
    *   **Body:**
        ```json
        {
          "vendor_id": 1,
          "item_description": "New Laptop",
          "quantity": 1,
          "estimated_price": 1500.00,
          "justification": "Developer machine upgrade"
        }
        ```
    *   **Response:** `201 Created` with the new requisition object.

*   **`GET /requisitions/my`**: Returns a list of PRs created by the logged-in user.
*   **`PUT /requisitions/{id}`**: Updates a requisition (if status is "Pending" and user is the requester).
*   **`DELETE /requisitions/{id}`**: Deletes a requisition (if status is "Pending" and user is the requester).
*   **`GET /requisitions/pending`** (Admin Only): Returns all PRs with "Pending" status.
*   **`GET /requisitions/all`** (Admin Only): Returns a list of all requisitions.
*   **`POST /requisitions/{id}/approve`** (Admin Only): Approves a PR and creates a Purchase Order.
*   **`POST /requisitions/{id}/reject`** (Admin Only): Rejects a PR.
*   **`PUT /admin/requisitions/{id}`** (Admin Only): Updates any requisition's details.
*   **`DELETE /admin/requisitions/{id}`** (Admin Only): Deletes any requisition.

### Purchase Orders

*All purchase order routes require authentication.*

*   **`GET /purchase-orders/all`** (Admin Only): Returns a list of all purchase orders.
*   **`GET /purchase-orders/{id}`**: Returns a purchase order by ID.
*   **`GET /purchase-orders/{id}/pdf`**: Generates and returns a PDF of the purchase order.

### Activity Log

*All activity log routes require authentication.*

*   **`GET /activity-logs`**: Returns a list of the last 100 activity events in the system.

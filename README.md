# Procurement System

This is a comprehensive Procurement Management System built from scratch, designed for enterprise use. It features a secure Go backend, a responsive Flutter frontend, and a PostgreSQL database.

## 🌟 Features

- **Full CRUD for Core Modules:** Admins can manage Users and Vendors. Employees can manage their own Purchase Requisitions.
- **Role-Based Access Control (RBAC):** Simple but effective roles (Admin, Employee, etc.) to control access to features.
- **Secure Authentication:** JWT-based authentication with password hashing (bcrypt).
- **Purchase Requisition Workflow:** Create, update, and delete requisitions. Admins can approve or reject them.
- **Automatic Purchase Order (PO) Generation:** POs are automatically created when a requisition is approved.
- **PDF Generation:** Download a PDF representation of a Purchase Order.
- **RESTful API:** A well-structured and documented RESTful API.

## 🛠️ Prerequisites

Before you begin, ensure you have the following installed:

- **Go:** Version 1.21 or higher.
- **Flutter:** Version 3.10 or higher.
- **PostgreSQL:** A running instance of PostgreSQL.
- **Docker:** (Optional) For containerized setup.

## 🚀 Installation & Environment Setup

To get the full system running, you need to set up the database, backend, and frontend.

1.  **Database Setup:**
    - Make sure your PostgreSQL server is running.
    - Create a dedicated database and user for this project.
    - The required schema will be applied by the backend service upon startup using migration files located in `backend/migrations`.

2.  **Backend Setup:**
    - Navigate to the `backend` directory.
    - Follow the detailed instructions in `backend/README.md`.

3.  **Frontend Setup:**
    - Navigate to the `frontend` directory.
    - Follow the detailed instructions in `frontend/README.md`.

## 📖 API Documentation

The backend provides a RESTful API. For detailed information on endpoints, request/response formats, and authentication requirements, please refer to the API documentation within the `backend/README.md` file. We plan to add Swagger/OpenAPI documentation in the future.

## ☁️ Deployment Guidelines

(To be added) - This section will contain instructions for deploying the application to production environments. It will cover topics like building for production, environment configuration, and containerization with Docker.

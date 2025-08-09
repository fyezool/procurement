# Procurement System - Frontend

This directory contains the Flutter web frontend for the Procurement System.

## Prerequisites

- You need to have the Flutter SDK installed. If you don't have it, you can follow the instructions on the [official Flutter website](https://flutter.dev/docs/get-started/install).
- The backend server must be running. Please follow the instructions in the `backend/README.md` to start it. The frontend expects the backend to be available at `http://localhost:8080`.

## ✨ Features

The frontend application provides a user-friendly interface for interacting with the procurement system.

- **Login:** Secure login page to access the system.
- **Dynamic Navigation:** The side navigation menu is dynamically generated based on the user's role, ensuring users only see what they're supposed to.
- **User Management (Admin):** A data table view for listing all users, with options to edit their role/name or delete them. A form is available for creating new users.
- **Vendor Management (Admin):** A data table view for listing all vendors, with options to edit and delete them. A form is available for creating new vendors.
- **Requisition Management:**
  - A central hub for requisition-related tasks.
  - Employees can view a list of their own requisitions.
  - Employees can create, update, and delete their own requisitions (if they are in "Pending" status).
- **Approvals (Admin/Approver):** A dedicated screen to view all pending requisitions and approve or reject them with a single click.
- **Admin Overviews:** Admins have access to pages that list *all* requisitions and *all* purchase orders in the system.

## Setup & Running the App

1.  **Navigate to the frontend directory:**
    ```bash
    cd frontend
    ```

2.  **Install dependencies:**
    ```bash
    flutter pub get
    ```

3.  **Run the web app:**
    ```bash
    flutter run -d chrome
    ```
    This will start the application in a Chrome browser window.

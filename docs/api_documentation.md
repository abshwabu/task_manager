# Task Manager API Documentation

## Authentication

### Register

Creates a new user account. The first user to register will be an admin.

- **URL:** `/auth/register`
- **Method:** `POST`
- **Request Body:**

```json
{
  "username": "your_username",
  "password": "your_password"
}
```

- **Success Response:**

  - **Code:** 201 Created
  - **Content:**

  ```json
  {
    "message": "user created successfully"
  }
  ```

### Login

Authenticates a user and returns a JWT token.

- **URL:** `/auth/login`
- **Method:** `POST`
- **Request Body:**

```json
{
  "username": "your_username",
  "password": "your_password"
}
```

- **Success Response:**

  - **Code:** 200 OK
  - **Content:**

  ```json
  {
    "token": "your_jwt_token"
  }
  ```

## Tasks

All task endpoints require a valid JWT token in the `Authorization` header.

Example: `Authorization: Bearer your_jwt_token`

### Get All Tasks

- **URL:** `/api/tasks`
- **Method:** `GET`
- **Success Response:**

  - **Code:** 200 OK
  - **Content:** An array of task objects.

### Get Task by ID

- **URL:** `/api/tasks/:id`
- **Method:** `GET`
- **Success Response:**

  - **Code:** 200 OK
  - **Content:** A task object.

### Create Task (Admin only)

- **URL:** `/api/tasks`
- **Method:** `POST`
- **Request Body:**

```json
{
  "title": "Task Title",
  "description": "Task Description",
  "due_date": "2025-12-31T23:59:59Z",
  "status": "pending"
}
```

- **Success Response:**

  - **Code:** 201 Created
  - **Content:**

  ```json
  {
    "id": "task_id"
  }
  ```

### Update Task (Admin only)

- **URL:** `/api/tasks/:id`
- **Method:** `PUT`
- **Request Body:**

```json
{
  "title": "Updated Title",
  "description": "Updated Description",
  "due_date": "2026-01-15T23:59:59Z",
  "status": "in-progress"
}
```

- **Success Response:**

  - **Code:** 200 OK
  - **Content:**

  ```json
  {
    "message": "task updated successfully"
  }
  ```

### Delete Task (Admin only)

- **URL:** `/api/tasks/:id`
- **Method:** `DELETE`
- **Success Response:**

  - **Code:** 200 OK
  - **Content:**

  ```json
  {
    "message": "task deleted successfully"
  }
  ```

## Admin

### Promote User (Admin only)

- **URL:** `/api/admin/promote`
- **Method:** `POST`
- **Request Body:**

```json
{
  "username": "username_to_promote"
}
```

- **Success Response:**

  - **Code:** 200 OK
  - **Content:**

  ```json
  {
    "message": "user promoted successfully"
  }
  ```

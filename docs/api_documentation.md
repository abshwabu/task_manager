# Task Manager API Documentation

## Base URL
```
http://localhost:8080
```

## Prerequisites

### MongoDB Setup
This API requires MongoDB to be running locally.

**Installation:**
- Install MongoDB Community Edition
- Start MongoDB service: `mongod`
- Default connection: `mongodb://localhost:27017`
- Database: `taskmanager`
- Collection: `tasks`

**Dependencies:**
```bash
go mod tidy
```

## Endpoints

### 1. Get All Tasks
**GET** `/tasks`

Returns a list of all tasks.

**Response:**
- **Status Code:** 200 OK
- **Content-Type:** application/json

```json
[
  {
    "id": "1",
    "title": "Task 1",
    "description": "First task",
    "due_date": "2024-01-15",
    "status": "Pending"
  },
  {
    "id": "2",
    "title": "Task 2", 
    "description": "Second task",
    "due_date": "2024-01-16",
    "status": "In Progress"
  }
]
```

### 2. Get Task by ID
**GET** `/tasks/{id}`

Returns a specific task by its ID.

**Parameters:**
- `id` (path parameter) - Task ID

**Response:**
- **Status Code:** 200 OK (success) / 404 Not Found (task not found)
- **Content-Type:** application/json

**Success Response:**
```json
{
  "id": "1",
  "title": "Task 1",
  "description": "First task", 
  "due_date": "2024-01-15",
  "status": "Pending"
}
```

**Error Response:**
```json
{
  "message": "task not found"
}
```

### 3. Create New Task
**POST** `/tasks`

Creates a new task.

**Request Body:**
```json
{
  "id": "4",
  "title": "New Task",
  "description": "Task description",
  "due_date": "2024-01-20",
  "status": "Pending"
}
```

**Response:**
- **Status Code:** 201 Created
- **Content-Type:** application/json

```json
{
  "id": "4",
  "title": "New Task",
  "description": "Task description",
  "due_date": "2024-01-20", 
  "status": "Pending"
}
```

### 4. Update Task
**PUT** `/tasks/{id}`

Updates an existing task.

**Parameters:**
- `id` (path parameter) - Task ID

**Request Body:**
```json
{
  "title": "Updated Task Title",
  "description": "Updated description",
  "due_date": "2024-01-25",
  "status": "In Progress"
}
```

**Response:**
- **Status Code:** 200 OK (success) / 404 Not Found (task not found)
- **Content-Type:** application/json

**Success Response:**
```json
{
  "id": "1",
  "title": "Updated Task Title",
  "description": "Updated description",
  "due_date": "2024-01-25",
  "status": "In Progress"
}
```

**Error Response:**
```json
{
  "message": "task not found"
}
```

### 5. Delete Task
**DELETE** `/tasks/{id}`

Deletes a task by its ID.

**Parameters:**
- `id` (path parameter) - Task ID

**Response:**
- **Status Code:** 200 OK (success) / 404 Not Found (task not found)
- **Content-Type:** application/json

**Success Response:**
```json
{
  "message": "task deleted successfully"
}
```

**Error Response:**
```json
{
  "message": "task not found"
}
```

## Data Model

### Task Object
```json
{
  "id": "string",
  "title": "string", 
  "description": "string",
  "due_date": "string (ISO 8601 format)",
  "status": "string"
}
```

**Field Descriptions:**
- `id`: MongoDB ObjectID (24-character hex string)
- `title`: Task title
- `description`: Detailed description of the task
- `due_date`: Due date in ISO 8601 format
- `status`: Current status (e.g., "Pending", "In Progress", "Completed")

**Note:** Task IDs are automatically generated MongoDB ObjectIDs when creating new tasks.

## Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `404 Not Found`: Resource not found
- `400 Bad Request`: Invalid request body

## Example Usage

### cURL Examples

**Get all tasks:**
```bash
curl -X GET http://localhost:8080/tasks
```

**Get task by ID:**
```bash
curl -X GET http://localhost:8080/tasks/65a1b2c3d4e5f6789abcdef0
```

**Create new task:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "New Task",
    "description": "Task description", 
    "due_date": "2024-01-20T00:00:00Z",
    "status": "Pending"
  }'
```

**Update task:**
```bash
curl -X PUT http://localhost:8080/tasks/65a1b2c3d4e5f6789abcdef0 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Task",
    "description": "Updated description",
    "due_date": "2024-01-25T00:00:00Z", 
    "status": "In Progress"
  }'
```

**Delete task:**
```bash
curl -X DELETE http://localhost:8080/tasks/65a1b2c3d4e5f6789abcdef0
```
# Task Management API

This is a Task Management API built with Go, Gin, and MongoDB.

## Running Tests

To run the complete test suite, including unit, mock-based, and integration tests, ensure you have a MongoDB instance running locally on `mongodb://localhost:27017`.

To execute all tests and generate a coverage report, navigate to the project root directory and run the following command:

```bash
go test ./... -cover
```

This command will:
- Run all `_test.go` files in the project.
- Utilize Testify for assertions and mocking.
- Connect to a test MongoDB database named `task_manager_test` for repository integration tests.
- Print a summary of the test results and a coverage percentage.

To view a detailed HTML coverage report, you can run:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Structure

Tests are organized under the `Tests/` directory, following the clean architecture principles:

- `Tests/mocks/`: Contains Testify mocks for repository interfaces.
- `Tests/infrastructure/`: Unit tests for `jwt_service.go` and `password_service.go`.
- `Tests/usecases/`: Unit tests for `user_usecases.go` and `task_usecases.go` using mocked repositories.
- `Tests/middleware/`: Tests for `auth_middleware.go`.
- `Tests/controllers/`: Tests for `controller.go`.
- `Tests/routers/`: Tests for `router.go` endpoints using `httptest`.
- `Tests/repositories_integration/`: Integration tests for `task_repository.go` and `user_repository.go` against a live MongoDB instance.

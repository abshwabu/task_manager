package controllers_test

import (
	"bytes"
	domain "example/task_manager/Domain"
	controllers "example/task_manager/Delivery/controllers"
	infrastructure "example/task_manager/Infrastructure"
	mocks "example/task_manager/Tests/mocks"
	usecases "example/task_manager/Usecases"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupRouterAndController() (*gin.Engine, *controllers.Controller, *mocks.MockTaskRepository, *mocks.MockUserRepository) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockTaskRepo := new(mocks.MockTaskRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	taskUsecase := usecases.NewTaskUsecase(mockTaskRepo)
	userUsecase := usecases.NewUserUsecase(mockUserRepo)

	controller := controllers.NewController(taskUsecase, userUsecase)

	return r, controller, mockTaskRepo, mockUserRepo
}

func TestController_GetTasks(t *testing.T) {
	r, controller, mockTaskRepo, _ := setupRouterAndController()
	r.GET("/tasks", controller.GetTasks)

	expectedTasks := []domain.Task{
		{ID: primitive.NewObjectID(), Title: "Task 1"},
	}

	// Test case 1: Successful retrieval
	mockTaskRepo.On("GetAll").Return(expectedTasks, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var tasks []domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.NoError(t, err)
	assert.Equal(t, expectedTasks, tasks)
	mockTaskRepo.AssertExpectations(t)

	// Test case 2: Error from usecase
	mockTaskRepo.On("GetAll").Return(nil, errors.New("database error")).Once()

	req, _ = http.NewRequest(http.MethodGet, "/tasks", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database error")
	mockTaskRepo.AssertExpectations(t)
}

func TestController_GetTaskByID(t *testing.T) {
	r, controller, mockTaskRepo, _ := setupRouterAndController()
	r.GET("/tasks/:id", controller.GetTaskByID)

	taskID := primitive.NewObjectID().Hex()
	expectedTask := &domain.Task{ID: primitive.NewObjectID(), Title: "Test Task"}

	// Test case 1: Successful retrieval
	mockTaskRepo.On("GetByID", taskID).Return(expectedTask, nil).Once()

	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+taskID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var task domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &task)
	assert.NoError(t, err)
	assert.Equal(t, expectedTask, &task)
	mockTaskRepo.AssertExpectations(t)

	// Test case 2: Task not found
	mockTaskRepo.On("GetByID", "nonexistentID").Return(nil, errors.New("task not found")).Once()

	req, _ = http.NewRequest(http.MethodGet, "/tasks/nonexistentID", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "task not found")
	mockTaskRepo.AssertExpectations(t)
}

func TestController_CreateTask(t *testing.T) {
	r, controller, mockTaskRepo, _ := setupRouterAndController()
	r.POST("/tasks", controller.CreateTask)

	newTask := domain.Task{
		Title:       "New Task",
		Description: "Description for new task",
		DueDate:     time.Now(),
		Status:      "pending",
	}

	// Test case 1: Successful creation
	mockTaskRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(func(task domain.Task) *domain.Task {
		task.ID = primitive.NewObjectID()
		return &task
	}, nil).Once()

	jsonTask, _ := json.Marshal(newTask)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonTask))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTask domain.Task
	err := json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.NoError(t, err)
	assert.NotEqual(t, primitive.NilObjectID, createdTask.ID)
	assert.Equal(t, newTask.Title, createdTask.Title)
	mockTaskRepo.AssertExpectations(t)

	// Test case 2: Invalid JSON
	req, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")

	// Test case 3: Error from usecase
	mockTaskRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(nil, errors.New("database error")).Once()

	req, _ = http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonTask))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "database error")
	mockTaskRepo.AssertExpectations(t)
}

func TestController_UpdateTask(t *testing.T) {
	r, controller, mockTaskRepo, _ := setupRouterAndController()
	r.PUT("/tasks/:id", controller.UpdateTask)

	taskID := primitive.NewObjectID().Hex()
	updatedTask := domain.Task{
		Title:       "Updated Task",
		Description: "Updated description",
		DueDate:     time.Now(),
		Status:      "completed",
	}

	// Test case 1: Successful update
	mockTaskRepo.On("Update", taskID, mock.AnythingOfType("domain.Task")).Return(nil).Once()

	jsonTask, _ := json.Marshal(updatedTask)
	req, _ := http.NewRequest(http.MethodPut, "/tasks/"+taskID, bytes.NewBuffer(jsonTask))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "task updated successfully")
	mockTaskRepo.AssertExpectations(t)

	// Test case 2: Task not found
	mockTaskRepo.On("Update", "nonexistentID", mock.AnythingOfType("domain.Task")).Return(errors.New("task not found")).Once()

	req, _ = http.NewRequest(http.MethodPut, "/tasks/nonexistentID", bytes.NewBuffer(jsonTask))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "task not found")
	mockTaskRepo.AssertExpectations(t)

	// Test case 3: Invalid JSON
	req, _ = http.NewRequest(http.MethodPut, "/tasks/"+taskID, bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}

func TestController_DeleteTask(t *testing.T) {
	r, controller, mockTaskRepo, _ := setupRouterAndController()
	r.DELETE("/tasks/:id", controller.DeleteTask)

	taskID := primitive.NewObjectID().Hex()

	// Test case 1: Successful deletion
	mockTaskRepo.On("Delete", taskID).Return(nil).Once()

	req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+taskID, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "task deleted successfully")
	mockTaskRepo.AssertExpectations(t)

	// Test case 2: Task not found
	mockTaskRepo.On("Delete", "nonexistentID").Return(errors.New("task not found")).Once()

	req, _ = http.NewRequest(http.MethodDelete, "/tasks/nonexistentID", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "task not found")
	mockTaskRepo.AssertExpectations(t)
}

func TestController_Register(t *testing.T) {
	r, controller, _, mockUserRepo := setupRouterAndController()
	r.POST("/register", controller.Register)

	registerReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "newuser",
		Password: "newpassword",
	}

	// Test case 1: Successful registration
	mockUserRepo.On("GetByUsername", registerReq.Username).Return(nil, errors.New("not found")).Once()
	mockUserRepo.On("Create", mock.AnythingOfType("domain.User")).Return(&domain.User{Username: registerReq.Username, IsAdmin: true}, nil).Once()

	jsonReq, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "user created successfully")
	mockUserRepo.AssertExpectations(t)

	// Test case 2: User already exists
	mockUserRepo.On("GetByUsername", registerReq.Username).Return(&domain.User{}, nil).Once()

	req, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "user already exists")
	mockUserRepo.AssertExpectations(t)

	// Test case 3: Invalid JSON
	req, _ = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}

func TestController_Login(t *testing.T) {
	r, controller, _, mockUserRepo := setupRouterAndController()
	r.POST("/login", controller.Login)

	loginReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testuser",
		Password: "password123",
	}
	hashedPassword, _ := infrastructure.HashPassword(loginReq.Password)
	existingUser := &domain.User{Username: loginReq.Username, Password: hashedPassword, IsAdmin: false}

	// Test case 1: Successful login
	mockUserRepo.On("GetByUsername", loginReq.Username).Return(existingUser, nil).Once()

	jsonReq, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
	mockUserRepo.AssertExpectations(t)

	// Test case 2: Invalid credentials (user not found)
	mockUserRepo.On("GetByUsername", "nonexistent").Return(nil, errors.New("not found")).Once()
	invalidLoginReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "nonexistent",
		Password: "password123",
	}
	jsonInvalidReq, _ := json.Marshal(invalidLoginReq)
	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonInvalidReq))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
	mockUserRepo.AssertExpectations(t)

	// Test case 3: Invalid credentials (incorrect password)
	mockUserRepo.On("GetByUsername", loginReq.Username).Return(existingUser, nil).Once()
	incorrectPasswordReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testuser",
		Password: "wrongpassword",
	}
	jsonIncorrectPasswordReq, _ := json.Marshal(incorrectPasswordReq)
	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonIncorrectPasswordReq))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid credentials")
	mockUserRepo.AssertExpectations(t)

	// Test case 4: Invalid JSON
	req, _ = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}

func TestController_PromoteUser(t *testing.T) {
	r, controller, _, mockUserRepo := setupRouterAndController()
	r.POST("/promote", controller.PromoteUser)

	promoteReq := struct {
		Username string `json:"username"`
	}{
		Username: "promoteuser",
	}

	// Test case 1: Successful promotion
	mockUserRepo.On("UpdateRole", promoteReq.Username, true).Return(nil).Once()

	jsonReq, _ := json.Marshal(promoteReq)
	req, _ := http.NewRequest(http.MethodPost, "/promote", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user promoted successfully")
	mockUserRepo.AssertExpectations(t)

	// Test case 2: User not found
	mockUserRepo.On("UpdateRole", "nonexistent", true).Return(errors.New("user not found")).Once()
	invalidPromoteReq := struct {
		Username string `json:"username"`
	}{
		Username: "nonexistent",
	}
	jsonInvalidReq, _ := json.Marshal(invalidPromoteReq)
	req, _ = http.NewRequest(http.MethodPost, "/promote", bytes.NewBuffer(jsonInvalidReq))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "user not found")
	mockUserRepo.AssertExpectations(t)

	// Test case 3: Invalid JSON
	req, _ = http.NewRequest(http.MethodPost, "/promote", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid character")
}

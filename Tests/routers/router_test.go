package routers_test

import (
	"bytes"
	domain "example/task_manager/Domain"
	controllers "example/task_manager/Delivery/controllers"
	"example/task_manager/Delivery/routers"
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

func setupTestEnvironment() (*gin.Engine, *mocks.MockTaskRepository, *mocks.MockUserRepository) {
	gin.SetMode(gin.TestMode)

	mockTaskRepo := new(mocks.MockTaskRepository)
	mockUserRepo := new(mocks.MockUserRepository)

	taskUsecase := usecases.NewTaskUsecase(mockTaskRepo)
	userUsecase := usecases.NewUserUsecase(mockUserRepo)

	controller := controllers.NewController(taskUsecase, userUsecase)
	router := routers.SetupRouter(controller)

	return router, mockTaskRepo, mockUserRepo
}

func getTokenForUser(username, password string, router *gin.Engine, mockUserRepo *mocks.MockUserRepository) (string, error) {
	hashedPassword, _ := infrastructure.HashPassword(password)
	existingUser := &domain.User{Username: username, Password: hashedPassword, IsAdmin: false}
	mockUserRepo.On("GetByUsername", username).Return(existingUser, nil).Once()

	loginReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	}
	jsonReq, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return "", errors.New("failed to get token")
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["token"], nil
}

func getAdminToken(username, password string, router *gin.Engine, mockUserRepo *mocks.MockUserRepository) (string, error) {
	hashedPassword, _ := infrastructure.HashPassword(password)
	adminUser := &domain.User{Username: username, Password: hashedPassword, IsAdmin: true}
	mockUserRepo.On("GetByUsername", username).Return(adminUser, nil).Once()

	loginReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	}
	jsonReq, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		return "", errors.New("failed to get admin token")
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp["token"], nil
}

func TestRouter_AuthRoutes(t *testing.T) {
	router, _, mockUserRepo := setupTestEnvironment()

	// Test /auth/register - success
	registerReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "newuser",
		Password: "newpassword",
	}
	mockUserRepo.On("GetByUsername", registerReq.Username).Return(nil, errors.New("not found")).Once()
	mockUserRepo.On("Create", mock.AnythingOfType("domain.User")).Return(&domain.User{Username: registerReq.Username, IsAdmin: true}, nil).Once()
	jsonReq, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonReq))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "user created successfully")
	mockUserRepo.AssertExpectations(t)

	// Test /auth/login - success
	loginReq := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "testuser",
		Password: "password123",
	}
	hashedPassword, _ := infrastructure.HashPassword(loginReq.Password)
	existingUser := &domain.User{Username: loginReq.Username, Password: hashedPassword, IsAdmin: false}
	mockUserRepo.On("GetByUsername", loginReq.Username).Return(existingUser, nil).Once()
	jsonLoginReq, _ := json.Marshal(loginReq)
	req, _ = http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(jsonLoginReq))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "token")
	mockUserRepo.AssertExpectations(t)
}

func TestRouter_ProtectedRoutes(t *testing.T) {
	router, mockTaskRepo, mockUserRepo := setupTestEnvironment()

	// Get a valid token for a regular user
	userToken, err := getTokenForUser("regularuser", "password", router, mockUserRepo)
	assert.NoError(t, err)
	assert.NotEmpty(t, userToken)

	// Test /api/tasks (GET) - success
	expectedTasks := []domain.Task{{ID: primitive.NewObjectID(), Title: "Protected Task"}}
	mockTaskRepo.On("GetAll").Return(expectedTasks, nil).Once()
	req, _ := http.NewRequest(http.MethodGet, "/api/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var tasks []domain.Task
	json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.Equal(t, expectedTasks, tasks)
	mockTaskRepo.AssertExpectations(t)

	// Test /api/tasks/:id (GET) - success
	taskID := primitive.NewObjectID().Hex()
	expectedTask := &domain.Task{ID: primitive.NewObjectID(), Title: "Specific Protected Task"}
	mockTaskRepo.On("GetByID", taskID).Return(expectedTask, nil).Once()
	req, _ = http.NewRequest(http.MethodGet, "/api/tasks/"+taskID, nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var task domain.Task
	json.Unmarshal(w.Body.Bytes(), &task)
	assert.Equal(t, expectedTask, &task)
	mockTaskRepo.AssertExpectations(t)

	// Test /api/tasks (GET) - unauthorized (no token)
	req, _ = http.NewRequest(http.MethodGet, "/api/tasks", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Authorization header required")
}

func TestRouter_AdminRoutes(t *testing.T) {
	router, mockTaskRepo, mockUserRepo := setupTestEnvironment()

	// Get a valid token for an admin user
	adminToken, err := getAdminToken("adminuser", "adminpassword", router, mockUserRepo)
	assert.NoError(t, err)
	assert.NotEmpty(t, adminToken)

	// Get a valid token for a regular user (for forbidden tests)
	userToken, err := getTokenForUser("regularuser_for_admin_test", "password", router, mockUserRepo)
	assert.NoError(t, err)
	assert.NotEmpty(t, userToken)

	// Test /api/tasks (POST) - admin success
	newTask := domain.Task{Title: "Admin Created Task", Description: "Desc", DueDate: time.Now(), Status: "pending"}
	mockTaskRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(func(task domain.Task) *domain.Task {
		task.ID = primitive.NewObjectID()
		return &task
	}, nil).Once()
	jsonNewTask, _ := json.Marshal(newTask)
	req, _ := http.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(jsonNewTask))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Admin Created Task")
	mockTaskRepo.AssertExpectations(t)

	// Test /api/tasks (POST) - user forbidden
	req, _ = http.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(jsonNewTask))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Admin access required")

	// Test /api/tasks/:id (PUT) - admin success
	taskID := primitive.NewObjectID().Hex()
	updatedTask := domain.Task{Title: "Updated Admin Task", Description: "Updated Desc", DueDate: time.Now(), Status: "completed"}
	mockTaskRepo.On("Update", taskID, mock.AnythingOfType("domain.Task")).Return(nil).Once()
	jsonUpdatedTask, _ := json.Marshal(updatedTask)
	req, _ = http.NewRequest(http.MethodPut, "/api/tasks/"+taskID, bytes.NewBuffer(jsonUpdatedTask))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "task updated successfully")
	mockTaskRepo.AssertExpectations(t)

	// Test /api/tasks/:id (PUT) - user forbidden
	req, _ = http.NewRequest(http.MethodPut, "/api/tasks/"+taskID, bytes.NewBuffer(jsonUpdatedTask))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Admin access required")

	// Test /api/tasks/:id (DELETE) - admin success
	mockTaskRepo.On("Delete", taskID).Return(nil).Once()
	req, _ = http.NewRequest(http.MethodDelete, "/api/tasks/"+taskID, nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "task deleted successfully")
	mockTaskRepo.AssertExpectations(t)

	// Test /api/tasks/:id (DELETE) - user forbidden
	req, _ = http.NewRequest(http.MethodDelete, "/api/tasks/"+taskID, nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Admin access required")

	// Test /api/admin/promote (POST) - admin success
	promoteReq := struct {
		Username string `json:"username"`
	}{
		Username: "user_to_promote",
	}
	mockUserRepo.On("UpdateRole", promoteReq.Username, true).Return(nil).Once()
	jsonPromoteReq, _ := json.Marshal(promoteReq)
	req, _ = http.NewRequest(http.MethodPost, "/api/admin/promote", bytes.NewBuffer(jsonPromoteReq))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user promoted successfully")
	mockUserRepo.AssertExpectations(t)

	// Test /api/admin/promote (POST) - user forbidden
	req, _ = http.NewRequest(http.MethodPost, "/api/admin/promote", bytes.NewBuffer(jsonPromoteReq))
	req.Header.Set("Authorization", "Bearer "+userToken)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Admin access required")
}

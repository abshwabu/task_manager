package usecases_test

import (
	domain "example/task_manager/Domain"
	mocks "example/task_manager/Tests/mocks"
	usecases "example/task_manager/Usecases"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

func TestTaskUsecase_GetAllTasks(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	expectedTasks := []domain.Task{
		{ID: primitive.NewObjectID(), Title: "Task 1"},
		{ID: primitive.NewObjectID(), Title: "Task 2"},
	}

	// Test case 1: Successful retrieval
	mockRepo.On("GetAll").Return(expectedTasks, nil).Once()

	tasks, err := taskUsecase.GetAllTasks()
	assert.NoError(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, expectedTasks, tasks)
	mockRepo.AssertExpectations(t)

	// Test case 2: Error from repository
	mockRepo.On("GetAll").Return(nil, errors.New("database error")).Once()

	tasks, err = taskUsecase.GetAllTasks()
	assert.Error(t, err)
	assert.Nil(t, tasks)
	assert.EqualError(t, err, "database error")
	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_GetTaskByID(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	taskID := primitive.NewObjectID().Hex()
	expectedTask := &domain.Task{ID: primitive.NewObjectID(), Title: "Test Task"}

	// Test case 1: Successful retrieval
	mockRepo.On("GetByID", taskID).Return(expectedTask, nil).Once()

	task, err := taskUsecase.GetTaskByID(taskID)
	assert.NoError(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, expectedTask, task)
	mockRepo.AssertExpectations(t)

	// Test case 2: Task not found
	mockRepo.On("GetByID", "nonexistentID").Return(nil, errors.New("task not found")).Once()

	task, err = taskUsecase.GetTaskByID("nonexistentID")
	assert.Error(t, err)
	assert.Nil(t, task)
	assert.EqualError(t, err, "task not found")
	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_CreateTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	newTask := domain.Task{
		Title:       "New Task",
		Description: "Description for new task",
		DueDate:     time.Now(),
		Status:      "pending",
	}

	// Test case 1: Successful creation
	mockRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(func(task domain.Task) *domain.Task {
		task.ID = primitive.NewObjectID()
		return &task
	}, nil).Once()

	createdTask, err := taskUsecase.CreateTask(newTask)
	assert.NoError(t, err)
	assert.NotNil(t, createdTask)
	assert.NotEqual(t, primitive.NilObjectID, createdTask.ID)
	assert.Equal(t, newTask.Title, createdTask.Title)
	mockRepo.AssertExpectations(t)

	// Test case 2: Error from repository
	mockRepo.On("Create", mock.AnythingOfType("domain.Task")).Return(nil, errors.New("database error")).Once()

	createdTask, err = taskUsecase.CreateTask(newTask)
	assert.Error(t, err)
	assert.Nil(t, createdTask)
	assert.EqualError(t, err, "database error")
	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_UpdateTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	taskID := primitive.NewObjectID().Hex()
	updatedTask := domain.Task{
		Title:       "Updated Task",
		Description: "Updated description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      "completed",
	}

	// Test case 1: Successful update
	mockRepo.On("Update", taskID, updatedTask).Return(nil).Once()

	err := taskUsecase.UpdateTask(taskID, updatedTask)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test case 2: Task not found during update
	mockRepo.On("Update", "nonexistentID", updatedTask).Return(errors.New("task not found")).Once()

	err = taskUsecase.UpdateTask("nonexistentID", updatedTask)
	assert.Error(t, err)
	assert.EqualError(t, err, "task not found")
	mockRepo.AssertExpectations(t)
}

func TestTaskUsecase_DeleteTask(t *testing.T) {
	mockRepo := new(mocks.MockTaskRepository)
	taskUsecase := usecases.NewTaskUsecase(mockRepo)

	taskID := primitive.NewObjectID().Hex()

	// Test case 1: Successful deletion
	mockRepo.On("Delete", taskID).Return(nil).Once()

	err := taskUsecase.DeleteTask(taskID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test case 2: Task not found during deletion
	mockRepo.On("Delete", "nonexistentID").Return(errors.New("task not found")).Once()

	err = taskUsecase.DeleteTask("nonexistentID")
	assert.Error(t, err)
	assert.EqualError(t, err, "task not found")
	mockRepo.AssertExpectations(t)
}

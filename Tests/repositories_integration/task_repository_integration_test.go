package repositories_integration_test

import (
	"context"
	domain "example/task_manager/Domain"
	repositories "example/task_manager/Repositories"
	"log"
	"testing"
	"time"
	"errors"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ( // Declare variables at package level for broader use
	client *mongo.Client
	db     *mongo.Database
	tasksCollection *mongo.Collection
	taskRepo domain.TaskRepository
)

func TestMain(m *testing.M) {
	// Setup: Connect to MongoDB and initialize repository
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db = client.Database("task_manager_test") // Use a dedicated test database
	tasksCollection = db.Collection("tasks")
	taskRepo = repositories.NewTaskRepository(tasksCollection)

	// Run tests
	m.Run()

	// Teardown: Disconnect from MongoDB
	err = client.Disconnect(context.Background())
	if err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	}

	// Clean up the test database after all tests are done
	db.Drop(context.Background())
}

func setupTest(t *testing.T) {
	// Clear the tasks collection before each test
	_, err := tasksCollection.DeleteMany(context.Background(), bson.D{})
	assert.NoError(t, err, "Failed to clear tasks collection")
}

func TestTaskRepository_Create(t *testing.T) {
	setupTest(t)

	newTask := domain.Task{
		Title:       "Buy groceries",
		Description: "Milk, Eggs, Bread",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      "pending",
	}

	createdTask, err := taskRepo.Create(newTask)
	assert.NoError(t, err)
	assert.NotNil(t, createdTask)
	assert.NotEqual(t, primitive.NilObjectID, createdTask.ID)
	assert.Equal(t, newTask.Title, createdTask.Title)

	// Verify it's in the database
	var fetchedTask domain.Task
	err = tasksCollection.FindOne(context.Background(), bson.M{"_id": createdTask.ID}).Decode(&fetchedTask)
	assert.NoError(t, err)
	assert.Equal(t, createdTask.ID, fetchedTask.ID)
	assert.Equal(t, createdTask.Title, fetchedTask.Title)
}

func TestTaskRepository_GetAll(t *testing.T) {
	setupTest(t)

	task1 := domain.Task{Title: "Task 1", Status: "pending"}
	task2 := domain.Task{Title: "Task 2", Status: "completed"}
	_, err := tasksCollection.InsertMany(context.Background(), []interface{}{task1, task2})
	assert.NoError(t, err)

	tasks, err := taskRepo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, tasks, 2)
	// Since IDs are generated on insert, compare titles
	assert.Contains(t, []string{tasks[0].Title, tasks[1].Title}, task1.Title)
	assert.Contains(t, []string{tasks[0].Title, tasks[1].Title}, task2.Title)
}

func TestTaskRepository_GetByID(t *testing.T) {
	setupTest(t)

	existingTask := domain.Task{ID: primitive.NewObjectID(), Title: "Find Me", Status: "pending"}
	_, err := tasksCollection.InsertOne(context.Background(), existingTask)
	assert.NoError(t, err)

	// Test case 1: Found task
	foundTask, err := taskRepo.GetByID(existingTask.ID.Hex())
	assert.NoError(t, err)
	assert.NotNil(t, foundTask)
	assert.Equal(t, existingTask.ID, foundTask.ID)
	assert.Equal(t, existingTask.Title, foundTask.Title)

	// Test case 2: Task not found
	notFoundTask, err := taskRepo.GetByID(primitive.NewObjectID().Hex())
	assert.Error(t, err)
	assert.Nil(t, notFoundTask)
	assert.Contains(t, err.Error(), "task not found")

	// Test case 3: Invalid ID format
	invalidIDTask, err := taskRepo.GetByID("invalid-id")
	assert.Error(t, err)
	assert.Nil(t, invalidIDTask)
	assert.Contains(t, err.Error(), "invalid task ID")
}

func TestTaskRepository_Update(t *testing.T) {
	setupTest(t)

	existingTask := domain.Task{ID: primitive.NewObjectID(), Title: "Old Title", Description: "Old Desc", DueDate: time.Now(), Status: "pending"}
	_, err := tasksCollection.InsertOne(context.Background(), existingTask)
	assert.NoError(t, err)

	updatedTask := domain.Task{Title: "New Title", Description: "New Desc", DueDate: time.Now().Add(time.Hour), Status: "completed"}

	// Test case 1: Successful update
	err = taskRepo.Update(existingTask.ID.Hex(), updatedTask)
	assert.NoError(t, err)

	var fetchedTask domain.Task
	err = tasksCollection.FindOne(context.Background(), bson.M{"_id": existingTask.ID}).Decode(&fetchedTask)
	assert.NoError(t, err)
	assert.Equal(t, updatedTask.Title, fetchedTask.Title)
	assert.Equal(t, updatedTask.Description, fetchedTask.Description)
	assert.Equal(t, updatedTask.Status, fetchedTask.Status)

	// Test case 2: Task not found
	err = taskRepo.Update(primitive.NewObjectID().Hex(), updatedTask)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")

	// Test case 3: Invalid ID format
	err = taskRepo.Update("invalid-id", updatedTask)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid task ID")
}

func TestTaskRepository_Delete(t *testing.T) {
	setupTest(t)

	existingTask := domain.Task{ID: primitive.NewObjectID(), Title: "Delete Me", Status: "pending"}
	_, err := tasksCollection.InsertOne(context.Background(), existingTask)
	assert.NoError(t, err)

	// Test case 1: Successful deletion
	err = taskRepo.Delete(existingTask.ID.Hex())
	assert.NoError(t, err)

	// Verify it's deleted from the database
	var fetchedTask domain.Task
	err = tasksCollection.FindOne(context.Background(), bson.M{"_id": existingTask.ID}).Decode(&fetchedTask)
	assert.Error(t, err, "Task should not be found after deletion")
	assert.True(t, errors.Is(err, mongo.ErrNoDocuments))

	// Test case 2: Task not found
	err = taskRepo.Delete(primitive.NewObjectID().Hex())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")

	// Test case 3: Invalid ID format
	err = taskRepo.Delete("invalid-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid task ID")
}

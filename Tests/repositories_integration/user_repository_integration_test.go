package repositories_integration_test

import (
	"context"
	domain "example/task_manager/Domain"
	repositories "example/task_manager/Repositories"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ( // Declare variables at package level for broader use
	userRepo domain.UserRepository
)

func setupUserTest(t *testing.T) *mongo.Collection {
	// Setup: Connect to MongoDB and initialize repository for User tests
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	assert.NoError(t, err, "Failed to connect to MongoDB for user tests")

	userDb := client.Database("task_manager_test") // Use a dedicated test database
	usersCollection := userDb.Collection("users")
	userRepo = repositories.NewUserRepository(usersCollection)

	// Clear the users collection before each test
	_, err = usersCollection.DeleteMany(context.Background(), bson.D{})
	assert.NoError(t, err, "Failed to clear users collection")

	t.Cleanup(func() {
		// Teardown: Disconnect from MongoDB and drop collection after each test
		err = client.Disconnect(context.Background())
		assert.NoError(t, err, "Error disconnecting from MongoDB for user tests")
		// userDb.Drop(context.Background())
	})
	return usersCollection
}

func TestUserRepository_Create(t *testing.T) {
	usersCollection := setupUserTest(t)

	newUser := domain.User{
		Username: "testuser",
		Password: "hashedpassword",
		IsAdmin:  false,
	}

	createdUser, err := userRepo.Create(newUser)
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.NotEqual(t, primitive.NilObjectID, createdUser.ID)
	assert.Equal(t, newUser.Username, createdUser.Username)

	// Verify it's in the database
	var fetchedUser domain.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": createdUser.ID}).Decode(&fetchedUser)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID, fetchedUser.ID)
	assert.Equal(t, createdUser.Username, fetchedUser.Username)
}

func TestUserRepository_GetByUsername(t *testing.T) {
	usersCollection := setupUserTest(t)

	existingUser := domain.User{ID: primitive.NewObjectID(), Username: "findme", Password: "hash", IsAdmin: false}
	_, err := usersCollection.InsertOne(context.Background(), existingUser)
	assert.NoError(t, err)

	// Test case 1: Found user
	foundUser, err := userRepo.GetByUsername(existingUser.Username)
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, existingUser.ID, foundUser.ID)
	assert.Equal(t, existingUser.Username, foundUser.Username)

	// Test case 2: User not found
	notFoundUser, err := userRepo.GetByUsername("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, notFoundUser)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_UpdateRole(t *testing.T) {
	usersCollection := setupUserTest(t)

	existingUser := domain.User{ID: primitive.NewObjectID(), Username: "updaterole", Password: "hash", IsAdmin: false}
	_, err := usersCollection.InsertOne(context.Background(), existingUser)
	assert.NoError(t, err)

	// Test case 1: Promote user to admin
	err = userRepo.UpdateRole(existingUser.Username, true)
	assert.NoError(t, err)

	var fetchedUser domain.User
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": existingUser.ID}).Decode(&fetchedUser)
	assert.NoError(t, err)
	assert.True(t, fetchedUser.IsAdmin)

	// Test case 2: Demote user (should also work, though promote is the main use case)
	err = userRepo.UpdateRole(existingUser.Username, false)
	assert.NoError(t, err)
	err = usersCollection.FindOne(context.Background(), bson.M{"_id": existingUser.ID}).Decode(&fetchedUser)
	assert.NoError(t, err)
	assert.False(t, fetchedUser.IsAdmin)

	// Test case 3: User not found
	err = userRepo.UpdateRole("nonexistent", true)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

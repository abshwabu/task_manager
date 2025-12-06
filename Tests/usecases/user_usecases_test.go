package usecases_test

import (
	domain "example/task_manager/Domain"
	infrastructure "example/task_manager/Infrastructure"
	mocks "example/task_manager/Tests/mocks"
	usecases "example/task_manager/Usecases"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestUserUsecase_Register(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userUsecase := usecases.NewUserUsecase(mockRepo)

	username := "testuser"
	password := "password123"
	hashedPassword, _ := infrastructure.HashPassword(password)

	// Test case 1: Successful registration (first user, so admin)
	mockRepo.On("GetByUsername", username).Return(nil, errors.New("user not found")).Once()
	mockRepo.On("Create", mock.AnythingOfType("domain.User")).Return(func(user domain.User) *domain.User {
		user.ID = primitive.NewObjectID()
		user.Password = hashedPassword
		user.IsAdmin = true // Should be true for the first user
		return &user
	}, nil).Once()

	user, err := userUsecase.Register(username, password)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, username, user.Username)
	assert.True(t, user.IsAdmin)
	assert.True(t, infrastructure.CheckPasswordHash(password, user.Password))
	mockRepo.AssertExpectations(t)

	// Test case 2: User already exists
	existingUser := &domain.User{Username: username}
	mockRepo.On("GetByUsername", username).Return(existingUser, nil).Once()

	user, err = userUsecase.Register(username, password)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "user already exists")
	mockRepo.AssertExpectations(t)

	// Test case 3: Error from repository during user creation
	mockRepo.On("GetByUsername", "anotheruser").Return(nil, errors.New("user not found")).Once()
	mockRepo.On("Create", mock.AnythingOfType("domain.User")).Return(nil, errors.New("db create error")).Once()

	user, err = userUsecase.Register("anotheruser", password)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, "db create error")
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_Login(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userUsecase := usecases.NewUserUsecase(mockRepo)

	username := "testuser"
	password := "password123"
	hashedPassword, _ := infrastructure.HashPassword(password)

	existingUser := &domain.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: hashedPassword,
		IsAdmin:  false,
	}

	// Test case 1: Successful login
	mockRepo.On("GetByUsername", username).Return(existingUser, nil).Once()

	token, err := userUsecase.Login(username, password)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := infrastructure.ValidateJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, existingUser.IsAdmin, claims.IsAdmin)
	mockRepo.AssertExpectations(t)

	// Test case 2: User not found
	mockRepo.On("GetByUsername", "nonexistent").Return(nil, errors.New("user not found")).Once()

	token, err = userUsecase.Login("nonexistent", password)
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "invalid credentials")
	mockRepo.AssertExpectations(t)

	// Test case 3: Incorrect password
	mockRepo.On("GetByUsername", username).Return(existingUser, nil).Once()

	token, err = userUsecase.Login(username, "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, token)
	assert.EqualError(t, err, "invalid credentials")
	mockRepo.AssertExpectations(t)
}

func TestUserUsecase_PromoteUser(t *testing.T) {
	mockRepo := new(mocks.MockUserRepository)
	userUsecase := usecases.NewUserUsecase(mockRepo)

	username := "testuser"

	// Test case 1: Successful promotion
	mockRepo.On("UpdateRole", username, true).Return(nil).Once()

	err := userUsecase.PromoteUser(username)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)

	// Test case 2: User not found during promotion
	mockRepo.On("UpdateRole", "nonexistent", true).Return(errors.New("user not found")).Once()

	err = userUsecase.PromoteUser("nonexistent")
	assert.Error(t, err)
	assert.EqualError(t, err, "user not found")
	mockRepo.AssertExpectations(t)
}

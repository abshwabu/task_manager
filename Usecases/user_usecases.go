package usecases

import (
	"errors"
	domain "example/task_manager/Domain"
	infrastructure "example/task_manager/Infrastructure"
)

type UserUsecaseImpl struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &UserUsecaseImpl{userRepo: userRepo}
}

func (u *UserUsecaseImpl) Register(username, password string) (*domain.User, error) {
	_, err := u.userRepo.GetByUsername(username)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	hashedPassword, err := infrastructure.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// First user is admin
	count := 0 // In real implementation, check user count
	user := domain.User{
		Username: username,
		Password: hashedPassword,
		IsAdmin:  count == 0,
	}

	return u.userRepo.Create(user)
}

func (u *UserUsecaseImpl) Login(username, password string) (string, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !infrastructure.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	return infrastructure.GenerateJWT(user.Username, user.IsAdmin)
}

func (u *UserUsecaseImpl) PromoteUser(username string) error {
	return u.userRepo.UpdateRole(username, true)
}
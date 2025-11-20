package domain

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	DueDate     time.Time          `json:"due_date" bson:"due_date"`
	Status      string             `json:"status" bson:"status"`
}

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Password string             `json:"password" bson:"password"`
	IsAdmin  bool               `json:"is_admin" bson:"is_admin"`
}

type TaskRepository interface {
	GetAll() ([]Task, error)
	GetByID(id string) (*Task, error)
	Create(task Task) (*Task, error)
	Update(id string, task Task) error
	Delete(id string) error
}

type UserRepository interface {
	GetByUsername(username string) (*User, error)
	Create(user User) (*User, error)
	UpdateRole(username string, isAdmin bool) error
}

type TaskUsecase interface {
	GetAllTasks() ([]Task, error)
	GetTaskByID(id string) (*Task, error)
	CreateTask(task Task) (*Task, error)
	UpdateTask(id string, task Task) error
	DeleteTask(id string) error
}

type UserUsecase interface {
	Register(username, password string) (*User, error)
	Login(username, password string) (string, error)
	PromoteUser(username string) error
}
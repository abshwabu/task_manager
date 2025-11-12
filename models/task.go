package models

import (
	"time"


	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	DueDate     time.Time          `json:"due_date"`
	Status      string             `json:"status"`
	UserID      primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

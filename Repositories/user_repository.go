package repositories

import (
	"context"
	"errors"
	domain "example/task_manager/Domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryImpl struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) domain.UserRepository {
	return &UserRepositoryImpl{collection: collection}
}

func (r *UserRepositoryImpl) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.collection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (r *UserRepositoryImpl) Create(user domain.User) (*domain.User, error) {
	user.ID = primitive.NewObjectID()
	_, err := r.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateRole(username string, isAdmin bool) error {
	update := bson.M{"$set": bson.M{"is_admin": isAdmin}}
	result, err := r.collection.UpdateOne(context.TODO(), bson.M{"username": username}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}
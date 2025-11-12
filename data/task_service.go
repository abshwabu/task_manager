package data

import (
	"context"
	"errors"
	"task_manager/database"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var taskCollection *mongo.Collection

func InitTaskService() {
	taskCollection = database.GetCollection("taskmanager", "tasks")
}

func GetAllTasks() []models.Task {
	cursor, err := taskCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		return []models.Task{}
	}
	var tasks []models.Task
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		return []models.Task{}
	}
	return tasks
}

func GetTasksByUserID(userID primitive.ObjectID) []models.Task {
	cursor, err := taskCollection.Find(context.TODO(), bson.M{"userid": userID})
	if err != nil {
		return []models.Task{}
	}
	var tasks []models.Task
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		return []models.Task{}
	}
	return tasks
}

func GetTaskByID(id string) (*models.Task, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task ID")
	}
	var task models.Task
	err = taskCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&task)
	if err != nil {
		return nil, errors.New("task not found")
	}
	return &task, nil
}

func AddTask(newTask models.Task) (primitive.ObjectID, error) {
	newTask.ID = primitive.NewObjectID()
	result, err := taskCollection.InsertOne(context.TODO(), newTask)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func UpdateTask(id string, updated models.Task) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}
	update := bson.M{"$set": bson.M{"title": updated.Title, "description": updated.Description, "duedate": updated.DueDate, "status": updated.Status}}
	result, err := taskCollection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}

func DeleteTask(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task ID")
	}
	result, err := taskCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("task not found")
	}
	return nil
}
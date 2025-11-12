package data

import (
	"context"
	"errors"
	"log"
	"task_manager/database"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

func InitUserService() {
	userCollection = database.GetCollection("taskmanager", "users")
}

func CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	user.Password = ""

	// The first user to register is an admin
	count, err := userCollection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		return err
	}
	if count == 0 {
		user.Role = "admin"
	} else {
		user.Role = "user"
	}

	_, err = userCollection.InsertOne(context.TODO(), user)
	return err
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := userCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func PromoteUser(username string) error {
	_, err := userCollection.UpdateOne(
		context.TODO(),
		bson.M{"username": username},
		bson.M{"$set": bson.M{"role": "admin"}},
	)
	return err
}

func Login(username, password string) (*models.User, error) {
	user, err := GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	log.Printf("User: %+v\n", user)
	log.Printf("Hashed Password from DB: %s\n", user.HashedPassword)
	log.Printf("Password from request: %s\n", password)

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		log.Printf("Error comparing passwords: %v\n", err)
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
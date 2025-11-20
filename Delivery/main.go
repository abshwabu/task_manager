package main

import (
	"context"
	"example/task_manager/Delivery/controllers"
	"example/task_manager/Delivery/routers"
	repositories "example/task_manager/Repositories"
	usecases "example/task_manager/Usecases"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	// MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("taskmanager")
	taskCollection := db.Collection("tasks")
	userCollection := db.Collection("users")

	// Initialize repositories
	taskRepo := repositories.NewTaskRepository(taskCollection)
	userRepo := repositories.NewUserRepository(userCollection)

	// Initialize use cases
	taskUsecase := usecases.NewTaskUsecase(taskRepo)
	userUsecase := usecases.NewUserUsecase(userRepo)

	// Initialize controller
	controller := controllers.NewController(taskUsecase, userUsecase)

	// Setup router
	router := routers.SetupRouter(controller)

	// Start server
	log.Println("Server starting on :8080")
	router.Run(":8080")
}
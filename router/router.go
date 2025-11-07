package router

import (
	"example/task_manager/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/tasks", controllers.GetTasksController)
	router.GET("/tasks/:id", controllers.GetTaskByIDController)
	router.POST("/tasks", controllers.AddTaskController)
	router.PUT("/tasks/:id", controllers.UpdateTaskController)
	router.DELETE("/tasks/:id", controllers.DeleteTaskController)

	return router
}
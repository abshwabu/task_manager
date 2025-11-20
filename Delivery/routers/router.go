package routers

import (
	"example/task_manager/Delivery/controllers"
	infrastructure "example/task_manager/Infrastructure"
	"github.com/gin-gonic/gin"
)

func SetupRouter(controller *controllers.Controller) *gin.Engine {
	router := gin.Default()

	// Auth routes
	auth := router.Group("/auth")
	{
		auth.POST("/register", controller.Register)
		auth.POST("/login", controller.Login)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(infrastructure.AuthMiddleware())
	{
		// Task routes
		api.GET("/tasks", controller.GetTasks)
		api.GET("/tasks/:id", controller.GetTaskByID)
		
		// Admin only routes
		admin := api.Group("")
		admin.Use(infrastructure.AdminMiddleware())
		{
			admin.POST("/tasks", controller.CreateTask)
			admin.PUT("/tasks/:id", controller.UpdateTask)
			admin.DELETE("/tasks/:id", controller.DeleteTask)
			admin.POST("/admin/promote", controller.PromoteUser)
		}
	}

	return router
}
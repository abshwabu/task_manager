package router

import (
	"task_manager/controllers"
	"task_manager/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", controllers.GetTasks)
			tasks.GET("/:id", controllers.GetTask)
			tasks.POST("", middleware.AdminMiddleware(), controllers.CreateTask)
			tasks.PUT("/:id", middleware.AdminMiddleware(), controllers.UpdateTask)
			tasks.DELETE("/:id", middleware.AdminMiddleware(), controllers.DeleteTask)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/promote", controllers.Promote)
		}
	}

	return r
}

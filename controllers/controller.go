package controllers

import (
	"net/http"
	"task_manager/data"
	"task_manager/middleware"
	"task_manager/models"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := data.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	token, err := middleware.GenerateJWT(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
}

func Login(c *gin.Context) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := data.Login(creds.Username, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := middleware.GenerateJWT(user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Promote(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := data.PromoteUser(req.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to promote user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted successfully"})
}

func GetTasks(c *gin.Context) {
	role, _ := c.Get("role")
	if role.(string) == "admin" {
		c.JSON(http.StatusOK, data.GetAllTasks())
		return
	}

	username, _ := c.Get("username")
	user, err := data.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}
	c.JSON(http.StatusOK, data.GetTasksByUserID(user.ID))
}

func GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := data.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, _ := c.Get("username")
	user, err := data.GetUserByUsername(username.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}
	task.UserID = user.ID

	id, err := data.AddTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := data.UpdateTask(id, task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task updated successfully"})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := data.DeleteTask(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
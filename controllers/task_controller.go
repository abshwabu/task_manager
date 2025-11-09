package controllers

import (
	"example/task_manager/data"
	"example/task_manager/models"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetTasksController(c *gin.Context)  {
	c.IndentedJSON(http.StatusOK, data.GetAllTasks())
}
func GetTaskByIDController(c *gin.Context)  {
	id := c.Param("id")
	task, err := data.GetTaskByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, task)
}
func AddTaskController(c *gin.Context)  {
	var newTask models.Task
	if err := c.BindJSON(&newTask); err != nil {
		return
	}
	data.AddTask(newTask)
	c.IndentedJSON(http.StatusCreated, newTask)
}
func UpdateTaskController(c *gin.Context)  {
	id := c.Param("id")
	task, err := data.GetTaskByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	if err := c.BindJSON(&task); err != nil {
		return
	}
	data.UpdateTask(id,*task)
	c.IndentedJSON(http.StatusOK, task)
}

func DeleteTaskController(c *gin.Context) {
	id := c.Param("id")
	err := data.DeleteTask(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}


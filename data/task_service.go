package data

import (
	"example/task_manager/models"
	"sync"
	"time"
	"errors"
)

var (
	tasks = []models.Task {
		{ID: "1", Title: "Task 1", Description: "First task", DueDate: time.Now(), Status: "Pending"},
    {ID: "2", Title: "Task 2", Description: "Second task", DueDate: time.Now().AddDate(0, 0, 1), Status: "In Progress"},
    {ID: "3", Title: "Task 3", Description: "Third task", DueDate: time.Now().AddDate(0, 0, 2), Status: "Completed"},
	}
	mutex sync.Mutex
)

func GetAllTasks() []models.Task  {
	mutex.Lock()
	defer mutex.Unlock()
	return append([]models.Task(nil), tasks...)

}

func GetTaskByID(id string) (*models.Task, error)  {
	mutex.Lock()
	defer mutex.Unlock()

	for _, task := range tasks{
		if task.ID == id  {
			copy := task 
			return &copy, nil
		}
	}
	return nil, errors.New("task not found")
}

func AddTask(newTask models.Task) {
	mutex.Lock()
	defer mutex.Unlock()
	tasks = append(tasks, newTask)
}
func UpdateTask(id string,updated models.Task) error {
	mutex.Lock()
	defer mutex.Unlock()

	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Title = updated.Title
			tasks[i].Description = updated.Description
			return nil
		}
	}
	return errors.New("task not found")
}

func DeleteTask(id string) error {
	mutex.Lock()
	defer mutex.Unlock()
	for i,task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found")
}

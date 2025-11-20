package usecases

import (
	domain "example/task_manager/Domain"
)

type TaskUsecaseImpl struct {
	taskRepo domain.TaskRepository
}

func NewTaskUsecase(taskRepo domain.TaskRepository) domain.TaskUsecase {
	return &TaskUsecaseImpl{taskRepo: taskRepo}
}

func (u *TaskUsecaseImpl) GetAllTasks() ([]domain.Task, error) {
	return u.taskRepo.GetAll()
}

func (u *TaskUsecaseImpl) GetTaskByID(id string) (*domain.Task, error) {
	return u.taskRepo.GetByID(id)
}

func (u *TaskUsecaseImpl) CreateTask(task domain.Task) (*domain.Task, error) {
	return u.taskRepo.Create(task)
}

func (u *TaskUsecaseImpl) UpdateTask(id string, task domain.Task) error {
	return u.taskRepo.Update(id, task)
}

func (u *TaskUsecaseImpl) DeleteTask(id string) error {
	return u.taskRepo.Delete(id)
}
package service

import (
	"github.com/KAA295/lo/domain"
)

type TaskRepo interface {
	Get(id int) (domain.Task, error)
	Post(task domain.Task) int
	GetAll(statusFilter domain.Status) []domain.Task
}

type TaskService struct {
	repo TaskRepo
}

func NewTaskService(repo TaskRepo) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) Get(id int) (domain.Task, error) {
	return s.repo.Get(id)
}

func (s *TaskService) GetAll(statusFilter string) []domain.Task {
	return s.repo.GetAll(statusFilter)
}

func (s *TaskService) Post(task domain.Task) int {
	return s.repo.Post(task)
}

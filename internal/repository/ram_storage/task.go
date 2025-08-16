package ramstorage

import (
	"sync"

	"github.com/KAA295/lo/domain"
)

type TaskRepo struct {
	db     map[int]domain.Task
	nextID int
	mutex  sync.RWMutex
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		db:     make(map[int]domain.Task),
		nextID: 1,
	}
}

func (r *TaskRepo) Get(id int) (domain.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	task, ok := r.db[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}

	return task, nil
}

func (r *TaskRepo) Post(task domain.Task) int {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	task.ID = r.nextID
	r.db[r.nextID] = task
	r.nextID++
	return task.ID
}

func (r *TaskRepo) GetAll(statusFilter domain.Status) []domain.Task {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var tasks []domain.Task
	for _, task := range r.db {
		if statusFilter == "" || task.Status == statusFilter {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

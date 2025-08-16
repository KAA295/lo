package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/KAA295/lo/api/rest/types"
	"github.com/KAA295/lo/domain"
)

type TaskService interface {
	Get(id int) (domain.Task, error)
	GetAll(statusFilter domain.Status) []domain.Task
	Post(task domain.Task) int
}

type Logger interface {
	Log(action string, message string)
}

type TaskHandler struct {
	service TaskService
	logger  Logger
}

func NewTaskHandler(service TaskService, logger Logger) *TaskHandler {
	return &TaskHandler{
		service: service,
		logger:  logger,
	}
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) *http.ServeMux {
	mux.HandleFunc("GET /tasks", h.GetAllTasks)
	mux.HandleFunc("GET /tasks/{id}", h.GetTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	return mux
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetRequest(r)
	if err != nil {
		h.logger.Log("ERROR", fmt.Sprintf("GetTaskByID %d: failed to create request: %v", req.ID, err))
		types.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	task, err := h.service.Get(req.ID)
	if err != nil {
		h.logger.Log("ERROR", fmt.Sprintf("GetTaskByID %d: failed: %v", req.ID, err))
		if errors.Is(err, domain.ErrNotFound) {
			types.SendError(w, err.Error(), http.StatusNotFound)
			return
		} else {
			types.SendError(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	resp := types.GetTaskResponse{
		ID:     task.ID,
		Data:   task.Data,
		Status: task.Status,
	}

	h.logger.Log("INFO", fmt.Sprintf("GetTaskByID %d: success", req.ID))

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Log("ERROR", fmt.Sprintf("GetTaskByID: failed to encode json: %v", err))
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateSetRequest(r)
	if err != nil {
		h.logger.Log("ERROR", fmt.Sprintf("CreateTask: failed to create request: %v", err))
		types.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := h.service.Post(domain.Task{
		Title:  req.Title,
		Data:   req.Data,
		Status: req.Status,
	})

	h.logger.Log("INFO", fmt.Sprintf("Task created: ID=%d", id))

	resp := types.CreateTaskResponse{ID: id}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Log("ERROR", fmt.Sprintf("CreateTask: failed to encode json: %v", err))
	}
}

func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetAllRequest(r)
	if err != nil {
		h.logger.Log("ERROR", fmt.Sprintf("GetAllTasks: failed to create request: %v", err))
		types.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks := h.service.GetAll(req.StatusFilter)

	var resp types.GetAllTasksResponse

	tasksSlice := make([]types.TaskEntry, len(tasks))
	resp.Tasks = tasksSlice

	for i, task := range tasks {
		resp.Tasks[i] = types.TaskEntry{
			ID:     task.ID,
			Data:   task.Data,
			Status: task.Status,
		}
	}

	h.logger.Log("INFO", "GetAllTasks: success")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.logger.Log("ERROR", fmt.Sprintf("GetAllTasks: failed to encode json: %v", err))
	}
}

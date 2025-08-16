package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KAA295/lo/domain"
)

type GetTaskRequest struct {
	ID int `json:"id"`
}

func CreateGetRequest(r *http.Request) (*GetTaskRequest, error) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, fmt.Errorf("can't convert idParam to int: %w", err)
	}
	return &GetTaskRequest{ID: id}, nil
}

type GetTaskResponse struct {
	ID     int           `json:"id"`
	Data   string        `json:"data"`
	Status domain.Status `json:"status"`
}

type CreateTaskRequest struct {
	Title  string        `json:"title"`
	Data   string        `json:"data"`
	Status domain.Status `json:"status"`
}

func CreateSetRequest(r *http.Request) (*CreateTaskRequest, error) {
	req := CreateTaskRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %w", err)
	}

	if !validateStatus(req.Status) {
		return nil, errors.New("unexpected status value")
	}

	if req.Title == "" {
		return nil, errors.New("title cannot be empty")
	}

	return &req, nil
}

type CreateTaskResponse struct {
	ID int `json:"id"`
}

type GetAllTasksRequest struct {
	StatusFilter domain.Status `json:"status"`
}

func CreateGetAllRequest(r *http.Request) (*GetAllTasksRequest, error) {
	q := r.URL.Query()

	var req GetAllTasksRequest

	status := q.Get("status")

	if !validateStatus(status) {
		return nil, errors.New("unexpected status value")
	}

	req.StatusFilter = status
	return &req, nil
}

type TaskEntry struct {
	ID     int           `json:"id"`
	Data   string        `json:"data"`
	Status domain.Status `json:"status"`
}

type GetAllTasksResponse struct {
	Tasks []TaskEntry `json:"tasks"`
}

func validateStatus(status string) bool {
	if status != domain.Ready &&
		status != domain.Done &&
		status != domain.InProgress &&
		status != "" {
		return false
	}
	return true
}

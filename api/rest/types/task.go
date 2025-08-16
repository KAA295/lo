package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/KAA295/lo/domain"
)

type GetRequest struct {
	ID int `json:"id"`
}

func CreateGetRequest(r *http.Request) (*GetRequest, error) {
	idParam := r.PathValue("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return nil, fmt.Errorf("can't convert idParam to int: %w", err)
	}
	return &GetRequest{ID: id}, nil
}

type GetResponse struct {
	ID     int           `json:"id"`
	Data   string        `json:"data"`
	Status domain.Status `json:"status"`
}

type SetRequest struct {
	Data   string        `json:"data"`
	Status domain.Status `json:"status"`
}

func CreateSetRequest(r *http.Request) (*SetRequest, error) {
	req := SetRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %w", err)
	}

	if !validateStatus(req.Status) {
		return nil, errors.New("unexpected status value")
	}

	return &req, nil
}

type SetResponse struct {
	ID int `json:"id"`
}

type GetAllRequest struct {
	StatusFilter domain.Status `json:"status"`
}

func CreateGetAllRequest(r *http.Request) (*GetAllRequest, error) {
	q := r.URL.Query()

	var req GetAllRequest

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

type GetAllResponse struct {
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

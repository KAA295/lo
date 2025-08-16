package domain

type Status = string

const (
	Done       Status = "done"
	InProgress Status = "in_progress"
	Ready      Status = "ready"
)

type Task struct {
	ID     int
	Title  string
	Data   string
	Status Status
}

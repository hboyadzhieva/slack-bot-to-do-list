package models

// Task - part of a to-do-list
type Task struct {
	Id          int     `json:"id"`
	Status      string  `json:"status"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	List        *List   `json:"list"`
	Asignee     *Person `json:"asignee"`
}

type TaskService interface {
	Task(id int) (*Task, error)
	TasksFromList(list *List) ([]*Task, error)
	CreateTask(t *Task) error
	DeleteTask(id int) error
	AsignTask(t *Task, p *Person) error
	SetStatus(t *Task, status string) error
}

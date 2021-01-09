package models

//To-do list consists of tasks
type List struct {
	Id      int      `json:"id"`
	Title   string   `json:"id"`
	Channel *Channel `json:"channel"`
}

type ListService interface {
	List(id int) (*List, error)
	Lists() ([]*List, error)
	AllListsInChannel(c *Channel) ([]*List, error)
	CreateList(l *List) error
	DeleteList(id int) error
	AddTask(t *Task) error
}

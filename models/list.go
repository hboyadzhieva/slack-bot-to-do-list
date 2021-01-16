package models

//To-do list consists of tasks
type List struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	ChannelId int    `json:"channel"`
}

type ListService interface {
	FindListById(id int) (*List, error)
	AllLists() ([]*List, error)
	AllListsInChannel(c *Channel) ([]*List, error)
	CreateList(l *List) error
	DeleteList(id int) error
	AddTask(t *Task) error
}

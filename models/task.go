package models

// Task - part of a to-do-list
type Task struct {
	Id          int    `json:"id"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ListId      int    `json:"list"`
	AsigneeId   int    `json:"asignee"`
}

package models

//To-do list consists of tasks
type List struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	ChannelId int    `json:"channel"`
}

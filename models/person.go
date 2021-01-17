package models

// Person represents a slack user
type Person struct {
	Id      int    `json:"id"`
	SlackId string `json:"slackId"`
}

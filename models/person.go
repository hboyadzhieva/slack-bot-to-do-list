package models

// Person represents a slack user
type Person struct {
	Id      int `json:"id"`
	SlackId int `json:"slackId"`
}

type PersonService interface {
	User(id int) (*Person, error)
	Users() ([]*Person, error)
	CreateUser(p *Person) error
	AllUsersInChannel(c *Channel) ([]*Person, error)
}

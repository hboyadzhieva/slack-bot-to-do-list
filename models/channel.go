package models

// Slack channel has users in the channel
// a to do list belongs to a single channel

type Channel struct {
	Id      int    `json:"id"`
	SlackId string `json:"slackId"`
}

type ChannelService interface {
	Channel(id int) (*Channel, error)
}

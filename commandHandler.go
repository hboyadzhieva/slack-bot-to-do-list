package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hboyadzhieva/slack-bot-to-do-list/models"
	"github.com/hboyadzhieva/slack-bot-to-do-list/mysql"
	"github.com/nlopes/slack"
	"strings"
)

func handleCreateList(params *slack.Msg) (string, error) {
	wrongFormatMsg := "command syntax: /tdcreatelist -t <title> -d <description>"
	words := strings.Split(params.Text, " ")
	var msg string
	if len(words) != 2 {
		return wrongFormatMsg, nil
	} else if words[0] != "-t" {
		return wrongFormatMsg, nil
	} else {
		var channel = &models.Channel{
			SlackId: params.Channel,
		}
		var list = &models.List{
			Title: words[1],
		}
		db, err := sql.Open("mysql", "root:pass@tcp(127.0.0.1:3306)/slack")

		if err != nil {
			return "Error opening DB connection", err
		}

		defer db.Close()
		listService := mysql.ListService{db}
		channelService := mysql.ChannelService{db}
		err = channelService.GetOrCreateChannel(channel)
		if err != nil {
			return "Error creating channel", err
		}
		list.ChannelId = channel.Id
		err = listService.CreateList(list)
		if err != nil {
			return "Error creating list: ", err
		}

	}
	return fmt.Sprintf("List created %v", msg), nil
}

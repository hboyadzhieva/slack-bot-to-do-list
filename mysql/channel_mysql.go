package mysql

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hboyadzhieva/slack-bot-to-do-list/models"
	"time"
)

//Implement ChannelService interface
type ChannelService struct {
	DB *sql.DB
}

func (cs *ChannelService) FindChannelById(id int) (*models.Channel, error) {
	var channel models.Channel
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	row := cs.DB.QueryRowContext(ctx, "SELECT ID, SLACK_ID FROM channel WHERE ID=?", id)
	err := row.Scan(&channel.Id, &channel.SlackId)
	if err != nil {
		return nil, err
	}
	return &channel, nil
}

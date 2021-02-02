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

func (cs *ChannelService) GetOrCreateChannel(c *models.Channel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var exists bool
	row := cs.DB.QueryRowContext(ctx, "SELECT EXISTS (SELECT ID, SLACK_ID FROM channel WHERE SLACK_ID=?)", c.SlackId)
	if err := row.Scan(&exists); err != nil {
		return err
	} else if exists {
		row := cs.DB.QueryRowContext(ctx, "SELECT ID, SLACK_ID FROM channel WHERE SLACK_ID=?", c.SlackId)
		if err := row.Scan(&c.Id, &c.SlackId); err != nil {
			return err
		}
	} else if !exists {
		query := "INSERT INTO channel (SLACK_ID) VALUES(?)"
		stmt, err := cs.DB.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()
		result, err := stmt.ExecContext(ctx, c.SlackId)
		if err != nil {
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		c.Id = int(id)
	}
	return nil
}

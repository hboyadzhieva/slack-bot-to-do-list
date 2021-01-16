package mysql

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hboyadzhieva/slack-bot-to-do-list/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

var channel = &models.Channel{
	Id:      1,
	SlackId: "SID12",
}

func TestFindChannelById(t *testing.T) {
	db, mock := NewMock()
	mockService := &ChannelService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, SLACK_ID FROM channel WHERE ID=\\?"
	rows := sqlmock.NewRows([]string{"ID", "SLACK_ID"}).
		AddRow(channel.Id, channel.SlackId)
	mock.ExpectQuery(query).WithArgs(channel.Id).WillReturnRows(rows)
	ch, err := mockService.FindChannelById(channel.Id)
	assert.NotNil(t, ch)
	assert.NoError(t, err)
}

func TestFindChannelByIdError(t *testing.T) {
	db, mock := NewMock()
	mockService := &ChannelService{db}
	defer func() {
		mockService.DB.Close()
	}()

	query := "SELECT ID, SLACK_ID FROM channel WHERE ID=\\?"
	rows := sqlmock.NewRows([]string{"ID", "SLACK_ID"})
	mock.ExpectQuery(query).WithArgs(channel.Id).WillReturnRows(rows)
	ch, err := mockService.FindChannelById(channel.Id)
	assert.Empty(t, ch)
	assert.Error(t, err)
}

package tododo

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hboyadzhieva/slack-bot-to-do-list/mysql"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandleCommandCallHelpHandler(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	mockRepo := &mysql.TaskRepository{db}
	mockHandler := CommandHandler{mockRepo}
	mockSlashCommand := &slack.SlashCommand{Command: "/tododo-help"}
	_, err = mockHandler.HandleCommand(mockSlashCommand)
	assert.NoError(t, err)
}

func TestHandleCommandErrorUnrecognizedCommand(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Errorf("Failed to open sqlmock database: Error %s", err)
	}
	mockRepo := &mysql.TaskRepository{db}
	mockHandler := CommandHandler{mockRepo}
	mockSlashCommand := &slack.SlashCommand{Command: "/tododo-NONO"}
	_, err = mockHandler.HandleCommand(mockSlashCommand)
	assert.Error(t, err)
}

func ExampleHandleCommand() {
	fmt.Println(1 + 2)
	//Output: 3
}

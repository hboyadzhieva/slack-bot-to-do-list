package tododo

import (
	"fmt"
	"github.com/hboyadzhieva/slack-bot-to-do-list/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockRepo struct{}

func (repo *MockRepo) PersistTask(t *mysql.Task) error {
	return nil
}

func (repo *MockRepo) GetTaskById(id int) (*mysql.Task, error) {
	return &mysql.Task{Id: 1, Status: mysql.StatusOpen, Title: "MockTitle", AsigneeId: "U1", ChannelId: "CH1"}, nil
}

func (repo *MockRepo) GetAllInChannel(channelId string) ([]*mysql.Task, error) {
	tasks := []*mysql.Task{&mysql.Task{Id: 1, Status: mysql.StatusOpen, Title: "MockTitle", AsigneeId: "U1", ChannelId: "CH1"}}
	return tasks, nil
}

func (repo *MockRepo) GetAllInChannelAssignedTo(channelId string, assigneeId string) ([]*mysql.Task, error) {
	tasks := []*mysql.Task{&mysql.Task{Id: 1, Status: mysql.StatusOpen, Title: "MockTitle", AsigneeId: "U1", ChannelId: "CH1"}}
	return tasks, nil
}

func (repo *MockRepo) AssignTaskTo(taskId int, assigneeId string) error {
	return nil
}

func (repo *MockRepo) SetStatus(taskId int, status string) error {
	return nil
}

func TestHandleHelpCommand(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleHelpCommand()
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, HelpHeader)
	assert.Contains(t, stringRes, HelpBlock1Text)
	assert.Contains(t, stringRes, HelpBlock2Text)
	assert.Contains(t, stringRes, HelpBlock3Text)
	assert.Contains(t, stringRes, HelpBlock4Text)
	assert.Contains(t, stringRes, HelpBlock5Text)
}

func TestHandleAddCommand(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleAddCommand("MockTitle", "CH1")
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, AddHeader)
	assert.Contains(t, stringRes, "Task added")
	assert.Contains(t, stringRes, "MockTitle")
}

func TestHandleShowCommand(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleShowCommand("", "CH1")
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, ShowHeader)
	assert.Contains(t, stringRes, "MockTitle")
}

func TestHandleAssignCommand(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleAssignCommand("1 U1")
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, UpdateHeader)
	assert.Contains(t, stringRes, "Assigned:")
	assert.Contains(t, stringRes, "MockTitle")
}

func TestHandleAssingCommandBadArgs(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleAssignCommand("1")
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, AssignBadArgsText)
}

func TestHandleProgressCommand(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleProgressCommand("1")
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, UpdateHeader)
	assert.Contains(t, stringRes, "Status:")
}

func TestHandleProgressCommandBadArgs(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleProgressCommand("1 one go")
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, AssignBadArgsText)
}

func TestHandleDoneCommand(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleDoneCommand("1")
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, UpdateHeader)
	assert.Contains(t, stringRes, "Status:")
}

func TestHandleDoneCommandBadArgs(t *testing.T) {
	mockHandler := &CommandHandler{&MockRepo{}}
	result, err := mockHandler.HandleDoneCommand("1 one go")
	stringRes := string(result)
	assert.NoError(t, err)
	assert.Contains(t, stringRes, AssignBadArgsText)
}

func ExampleHandleCommand() {
	fmt.Println(1 + 2)
	//Output: 3
}

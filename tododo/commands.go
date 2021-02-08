/*The package introduces command handlers to return proper response to the commands of ToDo bot*/
package tododo

import (
	"encoding/json"
	"fmt"
	"github.com/hboyadzhieva/slack-bot-to-do-list/mysql"
	"github.com/nlopes/slack"
	"strconv"
	"strings"
)

const (
	showMessage          = "ToDo"
	addMessage           = "ToDo: Added task"
	badArgsAssignMessage = "Bad arguments. Please enter /tododo-assign [taskId] [@user]"
)

type CommandHandlerInterface interface {
	HandleCommand(c *slack.SlashCommand) ([]byte, error)
	HandleHelpCommand(c *slack.SlashCommand) ([]byte, error)
	HandleAddCommand(c *slack.SlashCommand) ([]byte, error)
	HandleShowCommand(c *slack.SlashCommand) ([]byte, error)
	HandleAssignCommand(c *slack.SlashCommand) ([]byte, error)
}

type CommandHandler struct {
	Repository *mysql.TaskRepository
}

// Pass the command to the proper command handlers
func (handler *CommandHandler) HandleCommand(c *slack.SlashCommand) ([]byte, error) {
	switch c.Command {
	case "/tododo-help":
		return handler.HandleHelpCommand(c)
	case "/tododo-add":
		return handler.HandleAddCommand(c)
	case "/tododo-show":
		return handler.HandleShowCommand(c)
	case "/tododo-assign":
		return handler.HandleAssignCommand(c)
	default:
		return []byte("No such command"), fmt.Errorf("No such command %s", c.Command)
	}

	return []byte("Success"), nil
}

func (handler *CommandHandler) HandleHelpCommand(c *slack.SlashCommand) ([]byte, error) {
	header := NewHeaderBlock("Welcome! ToDo do can:")
	div := NewDividerBlock()
	block1 := NewSectionTextBlock("mrkdwn", "*/tododo add [task]*: add a task to your ToDo list")
	block2 := NewSectionTextBlock("mrkdwn", "*/tododo show*: show the tasks in your ToDo list")
	block3 := NewSectionTextBlock("mrkdwn", "*/tododo assign [taskId] [@user]*: assign a task to a user")
	resp := NewResponse(header, div, block1, block2, block3)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(byt))
	return byt, nil
}

func (handler *CommandHandler) HandleAddCommand(c *slack.SlashCommand) ([]byte, error) {
	task := mysql.NewTask(c.Text, c.ChannelID)
	// TO DO error checking
	err := handler.Repository.PersistTask(task)
	if err != nil {
		return []byte(""), err
	}
	return []byte("ToDo: Added task " + task.Title), nil
}

func (handler *CommandHandler) HandleShowCommand(c *slack.SlashCommand) ([]byte, error) {
	tasks, err := handler.Repository.GetAllInChannel(c.ChannelID)
	if err != nil {
		return []byte(""), err
	}
	//TO DO error checking
	header := NewHeaderBlock(showMessage)
	div := NewDividerBlock()
	blocks := make([]*Block, 0)
	for _, t := range tasks {
		idTitle := NewField("mrkdwn", "*"+strconv.Itoa(t.Id)+"*: "+t.Title)
		emoji := NewField("mrkdwn", getStatusEmoji(t.Status))
		assignee := NewField("mrkdwn", t.AsigneeId)
		status := NewField("mrkdwn", getStatusName(t.Status))
		block := NewSectionFieldsBlock(idTitle, emoji, assignee, status)
		blocks = append(blocks, block)
	}
	args := make([]*Block, 0)
	args = append(args, header)
	args = append(args, div)
	for _, b := range blocks {
		args = append(args, b)
	}
	resp := NewResponse(args...)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(byt))
	return byt, nil
}

func (handler *CommandHandler) HandleAssignCommand(c *slack.SlashCommand) ([]byte, error) {
	args := strings.Split(c.Text, " ")
	if len(args) != 2 {
		return []byte(badArgsAssignMessage), nil
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return []byte(badArgsAssignMessage), nil
	}
	err = handler.Repository.AssignTaskTo(id, args[1])
	if err != nil {
		return []byte(""), err
	}
	return []byte("Task assigned."), nil
}

func getStatusEmoji(status string) string {
	switch status {
	case mysql.StatusOpen:
		return ":question:"
	case mysql.StatusInProgress:
		return ":hourglass_flowing_sand:"
	case mysql.StatusDone:
		return ":white_check_mark:"
	default:
		return ""
	}
}

func getStatusName(status string) string {
	switch status {
	case mysql.StatusOpen:
		return "Open"
	case mysql.StatusInProgress:
		return "In progress"
	case mysql.StatusDone:
		return "Done"
	default:
		return ""
	}
}

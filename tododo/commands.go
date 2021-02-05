/*The package introduces command handlers to return proper response to the commands of ToDo bot*/
package tododo

import (
	"fmt"
	"github.com/hboyadzhieva/slack-bot-to-do-list/mysql"
	"github.com/nlopes/slack"
)

const (
	helpMessage = "Thanks for calling ToDo!\n" +
		"/tododo-add [task] : add a task and it will be available in your channel's list\n" +
		"/tododo-show : show the channel's ToDo list\n" +
		"/tododo-assign [taskId] [@user]: assign this task to a user\n"
	addMessage = "ToDo: Added task"
)

type CommandHandlerInterface interface {
	HandleCommand(c *slack.SlashCommand) (string, error)
	HandleHelpCommand(c *slack.SlashCommand) (string, error)
	HandleAddCommand(c *slack.SlashCommand) (string, error)
	HandleShowCommand(c *slack.SlashCommand) (string, error)
	HandleAssignCommand(c *slack.SlashCommand) (string, error)
}

type CommandHandler struct {
	Repository *mysql.TaskRepository
}

// Pass the command to the proper command handlers
func (handler *CommandHandler) HandleCommand(c *slack.SlashCommand) (string, error) {
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
		return "No such command", fmt.Errorf("No such command %s", c.Command)
	}

	return "Success", nil
}

func (handler *CommandHandler) HandleHelpCommand(c *slack.SlashCommand) (string, error) {
	return helpMessage, nil
}

func (handler *CommandHandler) HandleAddCommand(c *slack.SlashCommand) (string, error) {
	task := mysql.NewTask(c.Text, c.ChannelID)
	handler.Repository.PersistTask(task)
	return "ToDo: Added task " + task.Title, nil
}

func (handler *CommandHandler) HandleShowCommand(c *slack.SlashCommand) (string, error) {
	return "", nil
}

func (handler *CommandHandler) HandleAssignCommand(c *slack.SlashCommand) (string, error) {
	return "", nil
}

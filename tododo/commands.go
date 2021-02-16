/*Package tododo introduces command handlers to return proper response to the commands of ToDo bot.
Response body is structured in json format that conforms to Slack's Block Kit UI framework https://api.slack.com/block-kit in order to display intuituve and properly formatted response in Slack*/
package tododo

import (
	"encoding/json"
	"fmt"
	"github.com/hboyadzhieva/slack-bot-to-do-list/mysql"
	"github.com/nlopes/slack"
	"strconv"
	"strings"
)

// CommandHandlerInterface introduces functions to pass commands to the proper command handlers and return body of response to be forwarded and displayed in Slack.
type CommandHandlerInterface interface {
	HandleCommand(c *slack.SlashCommand) ([]byte, error)
	HandleHelpCommand() ([]byte, error)
	HandleAddCommand(text string, channelID string) ([]byte, error)
	HandleShowCommand(text string, channelID string) ([]byte, error)
	HandleAssignCommand(text string) ([]byte, error)
	HandleProgressCommand(text string) ([]byte, error)
	HandleDoneCommand(text string) ([]byte, error)
}

// CommandHandler implements CommandHandlerInterface
type CommandHandler struct {
	Repository mysql.TaskRepositoryInterface
}

// HandleCommand passes the command to the proper command handlers
func (handler *CommandHandler) HandleCommand(c *slack.SlashCommand) ([]byte, error) {
	switch c.Command {
	case "/tododo-help":
		return handler.HandleHelpCommand()
	case "/tododo-add":
		return handler.HandleAddCommand(c.Text, c.ChannelID)
	case "/tododo-show":
		return handler.HandleShowCommand(c.Text, c.ChannelID)
	case "/tododo-assign":
		return handler.HandleAssignCommand(c.Text)
	case "/tododo-start":
		return handler.HandleProgressCommand(c.Text)
	case "/tododo-done":
		return handler.HandleDoneCommand(c.Text)
	}
	return nil, fmt.Errorf("Can't handle command")
}

// HandleHelpCommand handles /tododo-help and returns proper response or error
func (handler *CommandHandler) HandleHelpCommand() ([]byte, error) {
	header := NewHeaderBlock(HelpHeader)
	div := NewDividerBlock()
	block1 := NewSectionTextBlock(MarkdownType, HelpBlock1Text)
	block2 := NewSectionTextBlock(MarkdownType, HelpBlock2Text)
	block3 := NewSectionTextBlock(MarkdownType, HelpBlock3Text)
	block4 := NewSectionTextBlock(MarkdownType, HelpBlock4Text)
	block5 := NewSectionTextBlock(MarkdownType, HelpBlock5Text)
	resp := NewResponse(header, div, block1, block2, block3, block4, block5)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

// HandleAddCommand handles /tododo-add and returns proper response or error
func (handler *CommandHandler) HandleAddCommand(text string, channelID string) ([]byte, error) {
	task := mysql.NewTask(text, channelID)
	err := handler.Repository.PersistTask(task)
	if err != nil {
		return nil, err
	}
	header := NewHeaderBlock(AddHeader)
	div := NewDividerBlock()
	block1 := NewSectionTextBlock(MarkdownType, "*Task added*: "+task.Title)
	resp := NewResponse(header, div, block1)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

// HandleShowCommand handles /tododo-show and returns proper response or error
func (handler *CommandHandler) HandleShowCommand(text string, channelID string) ([]byte, error) {
	tasks, err := handler.Repository.GetAllInChannel(channelID)
	if err != nil {
		return nil, err
	}
	header := NewHeaderBlock(ShowHeader)
	div := NewDividerBlock()
	blocks := make([]*Block, 0)
	for _, t := range tasks {
		idTitle := NewField(MarkdownType, "*"+strconv.Itoa(t.ID)+"*: "+t.Title)
		emoji := NewField(MarkdownType, getStatusEmoji(t.Status))
		assignee := NewField(MarkdownType, t.AsigneeID)
		status := NewField(MarkdownType, getStatusName(t.Status))
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
	return byt, nil
}

// HandleAssignCommand handles /tododo-assign and returns proper response or error.
func (handler *CommandHandler) HandleAssignCommand(text string) ([]byte, error) {
	header := NewHeaderBlock(UpdateHeader)
	div := NewDividerBlock()
	if !ValidateAssignCommandText(text) {
		errBlock := NewSectionTextBlock("plain_text", AssignBadArgsText)
		response := NewResponse(header, div, errBlock)
		byt, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return byt, nil
	}
	args := strings.Split(text, " ")
	id, _ := strconv.Atoi(args[0])
	err := handler.Repository.AssignTaskTo(id, args[1])
	if err == mysql.ErrNoRowOrMoreThanOne {
		errBlock := NewSectionTextBlock("plain_text", NoSuchTaskIDText)
		response := NewResponse(header, div, errBlock)
		byt, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return byt, nil
	} else if err != nil {
		return nil, err
	}
	task, err := handler.Repository.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	block1 := NewSectionTextBlock("mrkdwn", "Assigned: "+task.Title+" - "+task.AsigneeID)
	resp := NewResponse(header, div, block1)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

// HandleProgressCommand handles /tododo-start command and returns proper response or error.
func (handler *CommandHandler) HandleProgressCommand(text string) ([]byte, error) {
	header := NewHeaderBlock(UpdateHeader)
	div := NewDividerBlock()
	if !ValidateStatusText(text) {
		errBlock := NewSectionTextBlock("plain_text", ProgressBadArgsText)
		response := NewResponse(header, div, errBlock)
		byt, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return byt, nil
	}
	args := strings.Split(text, " ")
	id, _ := strconv.Atoi(args[0])
	err := handler.Repository.SetStatus(id, mysql.StatusInProgress)
	if err == mysql.ErrNoRowOrMoreThanOne {
		errBlock := NewSectionTextBlock("plain_text", NoSuchTaskIDText)
		response := NewResponse(header, div, errBlock)
		byt, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return byt, nil
	} else if err != nil {
		return nil, err
	}
	task, err := handler.Repository.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	block1 := NewSectionTextBlock(MarkdownType, "Status: "+task.Title+" - "+task.Status)
	resp := NewResponse(header, div, block1)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

// HandleDoneCommand handles /tododo-done command and returns proper response or error.
func (handler *CommandHandler) HandleDoneCommand(text string) ([]byte, error) {
	header := NewHeaderBlock(UpdateHeader)
	div := NewDividerBlock()
	if !ValidateStatusText(text) {
		errBlock := NewSectionTextBlock("plain_text", DoneBadArgsText)
		response := NewResponse(header, div, errBlock)
		byt, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return byt, nil
	}
	args := strings.Split(text, " ")
	id, _ := strconv.Atoi(args[0])
	err := handler.Repository.SetStatus(id, mysql.StatusDone)
	if err == mysql.ErrNoRowOrMoreThanOne {
		errBlock := NewSectionTextBlock("plain_text", NoSuchTaskIDText)
		response := NewResponse(header, div, errBlock)
		byt, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}
		return byt, nil
	} else if err != nil {
		return nil, err
	}
	task, err := handler.Repository.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	block1 := NewSectionTextBlock(MarkdownType, "Status: "+task.Title+" - "+task.Status)
	resp := NewResponse(header, div, block1)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

// ValidateAssignCommandText validates the args of /tododo-assign are exactly 2 - positive integer and a string represetation of assignee. Return true if the text is valid.
func ValidateAssignCommandText(text string) bool {
	args := strings.Split(text, " ")
	if len(args) != 2 {
		return false
	}
	num, err := strconv.Atoi(args[0])
	if err != nil || num < 1 {
		return false
	}
	return true
}

// ValidateStatusText validates the arg of /tododo-start and /tododo-done is exactly 1 - positive integer. Return true if the text is valid.
func ValidateStatusText(text string) bool {
	args := strings.Split(text, " ")
	if len(args) != 1 {
		return false
	}
	id, err := strconv.Atoi(args[0])
	if err != nil || id < 1 {
		return false
	}
	return true
}

func getStatusEmoji(status string) string {
	switch status {
	case mysql.StatusOpen:
		return StatusOpenEmoji
	case mysql.StatusInProgress:
		return StatusInProgressEmoji
	case mysql.StatusDone:
		return StatusDoneEmoji
	default:
		return ""
	}
}

func getStatusName(status string) string {
	switch status {
	case mysql.StatusOpen:
		return StatusOpenText
	case mysql.StatusInProgress:
		return StatusInProgressText
	case mysql.StatusDone:
		return StatusDoneText
	default:
		return ""
	}
}

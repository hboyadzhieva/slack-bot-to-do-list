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

type CommandHandlerInterface interface {
	HandleCommand(c *slack.SlashCommand) ([]byte, error)
	HandleHelpCommand() ([]byte, error)
	HandleAddCommand(text string, channelId string) ([]byte, error)
	HandleShowCommand(text string, channelId string) ([]byte, error)
	HandleAssignCommand(text string) ([]byte, error)
	HandleProgressCommand(text string) ([]byte, error)
	HandleDoneCommand(text string) ([]byte, error)
}

type CommandHandler struct {
	Repository mysql.TaskRepositoryInterface
}

// Pass the command to the proper command handlers
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
	default:
		return []byte("No such command"), fmt.Errorf("No such command %s", c.Command)
	}

	return []byte("Success"), nil
}

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
	//fmt.Println(string(byt))
	return byt, nil
}

func (handler *CommandHandler) HandleAddCommand(text string, channelId string) ([]byte, error) {
	task := mysql.NewTask(text, channelId)
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

func (handler *CommandHandler) HandleShowCommand(text string, channelId string) ([]byte, error) {
	tasks, err := handler.Repository.GetAllInChannel(channelId)
	if err != nil {
		return []byte(""), err
	}
	//TO DO error checking
	header := NewHeaderBlock(ShowHeader)
	div := NewDividerBlock()
	blocks := make([]*Block, 0)
	for _, t := range tasks {
		idTitle := NewField(MarkdownType, "*"+strconv.Itoa(t.Id)+"*: "+t.Title)
		emoji := NewField(MarkdownType, getStatusEmoji(t.Status))
		assignee := NewField(MarkdownType, t.AsigneeId)
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
	//fmt.Println(string(byt))
	return byt, nil
}

func (handler *CommandHandler) HandleAssignCommand(text string) ([]byte, error) {
	args := strings.Split(text, " ")
	if len(args) != 2 {
		return []byte(AssignBadArgsText), nil
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return []byte(AssignBadArgsText), nil
	}
	err = handler.Repository.AssignTaskTo(id, args[1])
	if err != nil {
		return []byte(""), err
	}
	task, err := handler.Repository.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	header := NewHeaderBlock(UpdateHeader)
	div := NewDividerBlock()
	block1 := NewSectionTextBlock("mrkdwn", "Assigned: "+task.Title+" - "+task.AsigneeId)
	resp := NewResponse(header, div, block1)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

func (handler *CommandHandler) HandleProgressCommand(text string) ([]byte, error) {
	args := strings.Split(text, " ")
	if len(args) != 1 {
		return []byte(AssignBadArgsText), nil
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return []byte(AssignBadArgsText), nil
	}
	err = handler.Repository.SetStatus(id, mysql.StatusInProgress)
	if err != nil {
		return nil, err
	}
	task, err := handler.Repository.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	header := NewHeaderBlock(UpdateHeader)
	div := NewDividerBlock()
	block1 := NewSectionTextBlock(MarkdownType, "Status: "+task.Title+" - "+task.Status)
	resp := NewResponse(header, div, block1)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return byt, nil
}

func (handler *CommandHandler) HandleDoneCommand(text string) ([]byte, error) {
	args := strings.Split(text, " ")
	if len(args) != 1 {
		return []byte(AssignBadArgsText), nil
	}
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return []byte(AssignBadArgsText), nil
	}
	err = handler.Repository.SetStatus(id, mysql.StatusDone)
	if err != nil {
		return nil, err
	}
	task, err := handler.Repository.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	header := NewHeaderBlock(UpdateHeader)
	div := NewDividerBlock()
	block1 := NewSectionTextBlock(MarkdownType, "Status: "+task.Title+" - "+task.Status)
	resp := NewResponse(header, div, block1)
	byt, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return byt, nil
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

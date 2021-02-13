package tododo

const (
	MarkdownType          = "mrkdwn"
	PlainTextType         = "plain_text"
	DividerType           = "divider"
	HelpHeader            = "Welcome! ToDo do can:"
	ShowHeader            = "ToDo"
	AddHeader             = "ToDo: Add task"
	UpdateHeader          = "ToDo: Task updated"
	AssignBadArgsText     = "Bad arguments. Please enter /tododo-assign [task ID] [@user]"
	HelpBlock1Text        = "*/tododo-add [task]*: add a task to your ToDo list"
	HelpBlock2Text        = "*/tododo-show*: show the tasks in your ToDo list"
	HelpBlock3Text        = "*/tododo-assign [taskId] [@user]*: assign a task to a user"
	HelpBlock4Text        = "*/tododo-start [taskId]*: start progress on a task"
	HelpBlock5Text        = "*/tododo-done [taskId]*: finish a task"
	StatusOpenEmoji       = ":question:"
	StatusInProgressEmoji = ":hourglass_flowing_sand:"
	StatusDoneEmoji       = ":white_check_mark:"
	StatusOpenText        = "Open"
	StatusInProgressText  = "In progress"
	StatusDoneText        = "Done"
)
